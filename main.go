package main

import (
	"context"
	"fmt"
	"log"

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

func main() {
	auditors, err := auditors.New(types.Config{})
	if err != nil {
		log.Fatalf("failed to initialize auditors: %s", err)
	}

	// Runners are executed in a specific order to resolve dependencies
	// correctly. Reordering this list may result in failed discovery.
	runners := []runner.Runner{
		kubelet.Runner,
		namespace.Runner,
		istiod.Runner,
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
}
