package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// old imports
	"context"
	"fmt"

	"github.com/fatih/color"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/debugz"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/runner/namespace"
	"github.com/praetorian-inc/mithril/pkg/types"
	"github.com/praetorian-inc/mithril/pkg/xds"

	_ "github.com/praetorian-inc/mithril/auditors/auth"
	_ "github.com/praetorian-inc/mithril/auditors/authz"
	_ "github.com/praetorian-inc/mithril/auditors/destinationrule"
	_ "github.com/praetorian-inc/mithril/auditors/gateway"
	_ "github.com/praetorian-inc/mithril/auditors/version"
)

var (
	// used for flags
	configFile string
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "mithril.yml", "mithril configuration file (default is `./mithril.yml`")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// we want to support config directories in pwd
		viper.AddConfigPath(".")

		// config file name is config.yaml
		viper.SetConfigName("mithril")
		viper.SetConfigType("yaml")
	}

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

	log.Infof("using config file: %s", viper.ConfigFileUsed())
}

// Execute is the entrypoint into the cmd line interface. It will execute the
// desired subcommand and check for an error, reporting it if so
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error during command execution: %s", err)
	}
}

func RunMithril() {
	auditors, err := auditors.New(types.Config{})
	if err != nil {
		log.Fatalf("failed to initialize auditors: %s", err)
	}

	runners := []runner.Runner{
		namespace.Runner,
	}
	var disco types.Discovery
	for _, r := range runners {
		log.Printf("running %s runner", r.Name)

		err := r.Run(&disco)
		if err != nil {
			log.Printf("failed to run %s: %s", r.Name, err)
		}
	}

	ctx := context.Background()
	var resources types.Resources
	if disco.DiscoveryAddress != "" {
		cli, err := xds.NewClient(disco.DiscoveryAddress)
		if err != nil {
			log.Printf("failed to initialize xds client: %s", err)
		}
		res, err := cli.Resources(ctx)
		if err != nil {
			log.Printf("failed to query xds resources: %s", err)
		}
		resources.Load(res)
		cli.Close()
	}
	if disco.DebugzAddress != "" {
		cli, err := debugz.NewClient(disco.DebugzAddress)
		if err != nil {
			log.Printf("failed to initialize debugz client: %s", err)
		}
		res, err := cli.Resources(ctx)
		if err != nil {
			log.Printf("failed to query debugz resources: %s", err)
		}
		resources.Load(res)
	}

	if resources.Len() == 0 {
		// TODO: import from static dir
		log.Fatalf("no resources discovered")
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
	for _, res := range results {
		fmt.Printf("%s: %s\n", red(res.Name), res.Description)
	}
}
