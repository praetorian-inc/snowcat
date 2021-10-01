package cli

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// old imports
	"github.com/fatih/color"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/runner/istiod"
	"github.com/praetorian-inc/mithril/pkg/runner/kubelet"
	"github.com/praetorian-inc/mithril/pkg/runner/namespace"
	"github.com/praetorian-inc/mithril/pkg/types"

    // Register all auditors
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
	inputDirectoryFlag   string
	exportDirectoryFlag  string
	istioVersionFlag     string
	istioNamespaceFlag   string
	discoveryAddressFlag string
	debugzAddressFlag    string
	kubeletAddressesFlag []string
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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("too many arguments specified")
		}
		if len(args) == 0 {
			return nil
		}

		path := args[0]
		stat, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("invalid input directory: %s", err)
		}
		if !stat.IsDir() {
			return fmt.Errorf("invalid input: %s must be a directory", path)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		RunMithril(args)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&configFileFlag, "config", "c", "mithril.yml", "mithril configuration file")
	rootCmd.Flags().StringVarP(&logLevelFlag, "log-level", "l", "info", "log level, see https://github.com/sirupsen/logrus#level-logging for options.")
	viper.BindPFlag("log-level", rootCmd.Flags().Lookup("log-level"))

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&inputDirectoryFlag, "input", "", "the input directory of static yaml files to scan")
	rootCmd.Flags().StringVar(&exportDirectoryFlag, "export", "", "write discovered resources to the specified export directory as yaml")

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

	rootCmd.Flags().BoolVarP(&saveConfFlag, "save-config", "s", false, "whether or not to save discovery to current config file")
}

func initConfig() {
	viper.SetConfigFile(configFileFlag)

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("configuration file not found")
		} else if os.IsNotExist(err) {
			log.Info("configuration file not found")
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
	}
}

func saveFinalDiscovery(disco types.Discovery) {
	viper.Set("istio-version", disco.IstioVersion)
	viper.Set("istio-namespace", disco.IstioNamespace)
	viper.Set("discovery-address", disco.DiscoveryAddress)
	viper.Set("debugz-address", disco.DebugzAddress)
	viper.Set("kubelet-addresses", disco.KubeletAddresses)
}

func RunMithril(args []string) {
	var inputPath string
	if len(args) == 1 {
		inputPath = args[0]
	}

	auditors, err := auditors.New(types.Config{})
	if err != nil {
		log.Fatalf("failed to initialize auditors: %s", err)
	}

	disco := buildInitialDiscovery()
	resources := types.NewResources()

	if inputPath == "" {
		// Runners are executed in a specific order to resolve dependencies
		// correctly. Reordering this list may result in failed discovery.
		runners := runner.Runners{
			kubelet.Runner,
			namespace.Runner,
			istiod.Runner,
		}
		runners.Run(&disco, &resources)
	} else {
		err = resources.LoadFromDirectory(inputPath)
		if err != nil {
			log.Fatalf("failed to load resources: %s", err)
		}
	}

	if resources.Len() == 0 {
		log.Fatalf("failed to discovery any resources")
	}

	if exportDirectoryFlag != "" {
		log.Printf("exporting resources to %s", exportDirectoryFlag)
		err = resources.Export(exportDirectoryFlag)
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

	if saveConfFlag {
		saveFinalDiscovery(disco)

		log.Info("saving configuration file based on new discoveries")

		err := viper.WriteConfig()
		if err != nil {
			log.Errorf("could not save configuration: %s", err)
		}
	}
}
