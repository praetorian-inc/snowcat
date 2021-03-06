// Copyright 2021 Praetorian Security, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/praetorian-inc/snowcat/auditors"
	// blank imports are for auditor registration
	_ "github.com/praetorian-inc/snowcat/auditors/authz"
	_ "github.com/praetorian-inc/snowcat/auditors/destinationrule"
	_ "github.com/praetorian-inc/snowcat/auditors/gateway"
	_ "github.com/praetorian-inc/snowcat/auditors/install"
	_ "github.com/praetorian-inc/snowcat/auditors/peerauth"
	_ "github.com/praetorian-inc/snowcat/auditors/version"
	"github.com/praetorian-inc/snowcat/pkg/runner"
	"github.com/praetorian-inc/snowcat/pkg/runner/istiod"
	"github.com/praetorian-inc/snowcat/pkg/runner/kubelet"
	"github.com/praetorian-inc/snowcat/pkg/runner/namespace"
	"github.com/praetorian-inc/snowcat/pkg/types"
)

var (
	configFileFlag       string
	logLevelFlag         string
	formatFlag           string
	exportDirectoryFlag  string
	outputFileFlag       string
	istioVersionFlag     string
	istioNamespaceFlag   string
	discoveryAddressFlag string
	debugzAddressFlag    string
	kubeletAddressesFlag []string
	saveConfFlag         bool
	jobMode              bool
)

const (
	jobCompleteMsg = `snowcat job complete! use the following command to export the results:

kubectl -n %s cp %s:%s snowcat-results
`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snowcat [input]",
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
		RunSnowcat(args)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&configFileFlag, "config", "c", "snowcat.yml",
		"snowcat configuration file")
	rootCmd.Flags().StringVarP(&logLevelFlag, "log-level", "l", "info",
		"log level, see https://github.com/sirupsen/logrus#level-logging for options.")
	viper.BindPFlag("log-level", rootCmd.Flags().Lookup("log-level"))

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&formatFlag, "format", "text", "output format [json, text]")

	rootCmd.Flags().StringVar(&exportDirectoryFlag, "export", "",
		"write discovered resources to the specified export directory as yaml")

	rootCmd.Flags().StringVar(&outputFileFlag, "output", "",
		"write results to the specified file")

	rootCmd.Flags().StringVar(&istioVersionFlag, "istio-version", "",
		"the version of the istio control plane")
	viper.BindPFlag("istio-version", rootCmd.Flags().Lookup("istio-version"))

	rootCmd.Flags().StringVar(&istioNamespaceFlag, "istio-namespace", "",
		"the kubernetes namespace of the istio control plane")
	viper.BindPFlag("istio-namespace", rootCmd.Flags().Lookup("istio-namespace"))

	rootCmd.Flags().StringVar(&discoveryAddressFlag, "discovery-address", "",
		"ip:port of istiod's unauthenticated xds")
	viper.BindPFlag("discovery-address", rootCmd.Flags().Lookup("discovery-address"))

	rootCmd.Flags().StringVar(&debugzAddressFlag, "debugz-address", "",
		"ip:port of istiod's debug api")
	viper.BindPFlag("debugz-address", rootCmd.Flags().Lookup("debugz-address"))

	rootCmd.Flags().StringSliceVar(&kubeletAddressesFlag, "kubelet-addresses", []string{},
		"list of addresses in form host:port of each node's kubelet read-only api")
	viper.BindPFlag("kubelet-addresses", rootCmd.Flags().Lookup("kubelet-addresses"))

	rootCmd.Flags().BoolVarP(&saveConfFlag, "save-config", "s", false,
		"whether or not to save discovery to current config file")

	rootCmd.Flags().BoolVarP(&jobMode, "job-mode", "", false,
		"used when running snowcat as a k8s job, delays exit to allow extracting results")
}

func initConfig() {
	viper.SetConfigFile(configFileFlag)

	// we would like for our env vars to have _'s but our vars to be -'d
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("configuration file not found")
		} else if os.IsNotExist(err) {
			log.Warn("configuration file not found")
		} else {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("issue reading configuration file")
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
		TimestampFormat: time.RFC3339,
	})

	log.WithFields(log.Fields{
		"configFile": viper.ConfigFileUsed(),
	}).Info("successfully loaded config file")
}

// Execute is the entrypoint into the cmd line interface. It will execute the
// desired subcommand and check for an error, reporting it if so
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("error during command execution")
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

// RunSnowcat runs the scanner.
func RunSnowcat(args []string) {
	var err error

	var inputPath string
	if len(args) == 1 {
		inputPath = args[0]
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
			log.WithFields(log.Fields{
				"err": err,
			}).Fatalf("failed to load resources")
		}
	}

	// TODO: generalize the empty disco check
	if resources.Len() == 0 && disco.IstioVersion == "" {
		log.Fatal("failed to discovery any resources")
	}

	if exportDirectoryFlag != "" {
		log.WithFields(log.Fields{
			"exportDirectory": exportDirectoryFlag,
		}).Info("exporting resources")

		err = resources.Export(exportDirectoryFlag)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Errorf("failed to export resources")
		}
	}

	var results []types.AuditResult
	for _, auditor := range auditors.All() {
		log.WithFields(log.Fields{
			"auditor": auditor.Name(),
		}).Info("running auditor")

		res, err := auditor.Audit(disco, resources)
		if err != nil {
			log.WithFields(log.Fields{
				"auditor": auditor.Name(),
				"err":     err,
			}).Error("auditor failed to run")
		}
		results = append(results, res...)
	}

	var out io.WriteCloser
	if outputFileFlag != "" {
		out, err = os.Create(outputFileFlag)
		defer out.Close()
	} else {
		out = os.Stdout
	}

	switch formatFlag {
	case "json":
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		_ = enc.Encode(results)
	case "text":
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		for _, res := range results {
			fmt.Fprintf(out, "%s [%s]: %s\n", red(res.Name), yellow(res.Resource), res.Description)
		}
	}

	if saveConfFlag {
		saveFinalDiscovery(disco)

		log.Info("saving configuration file based on new discoveries")

		err := viper.WriteConfig()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Errorf("could not save configuration")
		}
	}

	if jobMode && exportDirectoryFlag != "" {
		podName := os.Getenv("POD_NAME")
		if podName == "" {
			podName = "<pod>"
		}

		namespace := os.Getenv("POD_NAMESPACE")
		if namespace == "" {
			namespace = "<namespace>"
		}

		fmt.Printf(jobCompleteMsg, namespace, podName, exportDirectoryFlag)
		time.Sleep(5 * time.Minute)
	}
}
