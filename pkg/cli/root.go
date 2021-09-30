package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"

	// old imports
	"context"
	"fmt"

	"github.com/fatih/color"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/debugz"
	kubeletclient "github.com/praetorian-inc/mithril/pkg/kubelet"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/runner/istiod"
	"github.com/praetorian-inc/mithril/pkg/runner/kubelet"
	"github.com/praetorian-inc/mithril/pkg/runner/namespace"
	"github.com/praetorian-inc/mithril/pkg/types"
	"github.com/praetorian-inc/mithril/pkg/xds"

	_ "github.com/praetorian-inc/mithril/auditors/auth"
	_ "github.com/praetorian-inc/mithril/auditors/authz"
	_ "github.com/praetorian-inc/mithril/auditors/destinationrule"
	_ "github.com/praetorian-inc/mithril/auditors/gateway"
	_ "github.com/praetorian-inc/mithril/auditors/install"
	_ "github.com/praetorian-inc/mithril/auditors/version"
)

var (
	// used for flags
	configFileFlag       string
	logLevelFlag         string
	istioVersionFlag     string
	istioNamespaceFlag   string
	discoveryAddressFlag string
	debugzAddressFlag    string
	kubeletAddressesFlag []string
	istiodIPsFlag        []string
	saveConfFlag         bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mithril",
	Short: "an istio security scanner",
	Long: `this tool can be used by an organization looking to audit their own
istio service mesh, or by a security engineer looking to evaluate a customer's mesh.
it is capable of operating in a few different modes, including configuration files
and live clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		RunMithril()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&configFileFlag, "config", "c", "mithril.yml", "mithril configuration file")
	rootCmd.Flags().StringVarP(&logLevelFlag, "log-level", "l", "info", "log level, see https://github.com/sirupsen/logrus#level-logging for options.")
	viper.BindPFlag("log-level", rootCmd.Flags().Lookup("log-level"))

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&istioVersionFlag, "istio-version", "", "the version of the istio control plane")
	viper.BindPFlag("istio-version", rootCmd.Flags().Lookup("istio-version"))

	rootCmd.Flags().StringVar(&istioNamespaceFlag, "istio-namespace", "", "the kubernetes namespace of the istio control plane")
	viper.BindPFlag("istio-namespace", rootCmd.Flags().Lookup("istio-namespace"))

	rootCmd.Flags().StringVar(&discoveryAddressFlag, "discovery-address", "", "ip:port of istiod's unauthenticated xds")
	viper.BindPFlag("discovery-address", rootCmd.Flags().Lookup("discovery-address"))

	rootCmd.Flags().StringVar(&debugzAddressFlag, "debugz-address", "", "ip:port of istiod's debug api")
	viper.BindPFlag("debugz-address", rootCmd.Flags().Lookup("debugz-address"))

	rootCmd.Flags().StringSliceVar(&kubeletAddressesFlag, "kubelet-addresses", []string{}, "list of addresses in form host:port of each node's kubelet read-only api")
	viper.BindPFlag("kubelet-addresses", rootCmd.Flags().Lookup("kubelet-addresses"))

	rootCmd.Flags().StringSliceVar(&istiodIPsFlag, "istiod-ips", []string{}, "list of ip addresses that appear to be the istio control plane")
	viper.BindPFlag("istiod-ips", rootCmd.Flags().Lookup("istiod-ips"))

	rootCmd.Flags().BoolVarP(&saveConfFlag, "save-config", "s", false, "whether or not to save discovery to current config file")
	viper.BindPFlag("save-config", rootCmd.Flags().Lookup("save-config"))
}

func initConfig() {
	viper.SetConfigFile(configFileFlag)

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("configuration file not found")
		} else {
			log.Fatalf("issue reading configuration file: %s", err)
		}
	}

	level, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(level)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		// https://tools.ietf.org/html/rfc3339
		TimestampFormat: time.RFC3339Nano,
	})

	log.Infof("using config file: %s", viper.ConfigFileUsed())
}

// Execute is the entrypoint into the cmd line interface. It will execute the
// desired subcommand and check for an error, reporting it if so
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error during command execution: %s", err)
	}
}

func buildInitialDiscovery() types.Discovery {
	return types.Discovery{
		IstioVersion:     viper.GetString("istio-version"),
		IstioNamespace:   viper.GetString("istio-namespace"),
		DiscoveryAddress: viper.GetString("discovery-address"),
		DebugzAddress:    viper.GetString("debugz-address"),
		KubeletAddresses: viper.GetStringSlice("kubelet-addresses"),
		IstiodIPs:        viper.GetStringSlice("istiod-ips"),
	}
}

func RunMithril() {
	auditors, err := auditors.New(types.Config{})
	if err != nil {
		log.Fatalf("failed to initialize auditors: %s", err)
	}

	disco := buildInitialDiscovery()

	// Runners are executed in a specific order to resolve dependencies
	// correctly. Reordering this list may result in failed discovery.
	runners := []runner.Runner{
		kubelet.Runner,
		namespace.Runner,
		istiod.Runner,
	}
	for _, r := range runners {
		log.Printf("running %s runner", r.Name)

		err := r.Run(&disco)
		if err != nil {
			log.Printf("failed to run %s: %s", r.Name, err)
		}
	}

	ctx := context.Background()
	resources := types.NewResources()
	if disco.DiscoveryAddress != "" {
		log.Printf("querying xds at %s", disco.DiscoveryAddress)
		cli, err := xds.NewClient(disco.DiscoveryAddress)
		if err != nil {
			log.Printf("failed to initialize xds client: %s", err)
		}
		res, err := cli.Resources(ctx)
		if err != nil {
			log.Printf("failed to query xds resources: %s", err)
		}
		resources.Load(res)
		disco.IstioVersion, err = cli.Version(ctx)
		if err != nil {
			log.Printf("failed to query versions via xds resources: %s", err)
		}
		cli.Close()
	}
	if disco.DebugzAddress != "" {
		log.Printf("querying debug API at %s", disco.DebugzAddress)
		cli, err := debugz.NewClient(disco.DebugzAddress)
		if err != nil {
			log.Printf("failed to initialize debugz client: %s", err)
		}
		res, err := cli.Resources(ctx)
		if err != nil {
			log.Printf("failed to query debugz resources: %s", err)
		}
		disco.IstioVersion, err = cli.Version(ctx)
		if err != nil {
			log.Printf("failed to query versions via debugz resources: %s", err)
		}
		resources.Load(res)
	}
	if len(disco.KubeletAddresses) > 0 {
		for _, addr := range disco.KubeletAddresses {
			cli, err := kubeletclient.NewClient(addr)
			pods, err := cli.Pods(ctx)
			if err != nil {
				log.Printf("failed to query pods from kubelet: %s", err)
				continue
			}
			var res []runtime.Object
			for i := range pods {
				res = append(res, &pods[i])
			}
			resources.Load(res)
		}
	}

	if resources.Len() == 0 {
		// TODO: import from static dir
		err = resources.LoadFromDirectory("_fixtures")
		if err != nil {
			log.Fatalf("failed to load resources: %s", err)
		}
	} else {
		log.Printf("exporting resources")
		err = resources.Export("export")
		if err != nil {
			log.Printf("failed to export resources: %s", err)
		}
	}

	var results []types.AuditResult
	for _, auditor := range auditors {
		log.Printf("running auditor %s", auditor.Name())

		res, err := auditor.Audit(disco, resources)
		if err != nil {
			log.Printf("%s failed to run: %s", auditor.Name(), err)
		}
		results = append(results, res...)
	}

	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	for _, res := range results {
		fmt.Printf("%s [%s]: %s\n", red(res.Name), yellow(res.Resource), res.Description)
	}

	if save := viper.GetBool("save-config"); save {
		log.Info("saving configuration file based on new discoveries")

		err := viper.WriteConfig()
		if err != nil {
			log.Errorf("could not save configuration: %s", err)
		}
	}
}
