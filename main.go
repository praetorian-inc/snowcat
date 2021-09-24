package main

import (
	"fmt"
	"log"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/context/xds"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/runner/discovery"
	"github.com/praetorian-inc/mithril/pkg/runner/namespace"
	"github.com/praetorian-inc/mithril/pkg/types"

	_ "github.com/praetorian-inc/mithril/auditors/auth"
	_ "github.com/praetorian-inc/mithril/auditors/authz"
	_ "github.com/praetorian-inc/mithril/auditors/destinationrule"
	_ "github.com/praetorian-inc/mithril/auditors/gateway"
)

func main() {
	auditors, err := auditors.New(types.Config{})
	if err != nil {
		log.Fatalf("failed to initialize auditors: %s", err)
	}

	runners := []runner.Runner{
		namespace.Runner,
		discovery.Runner,
	}
	conf := make(map[string]string)
	for _, r := range runners {
		err := r.Run(conf)
		if err != nil {
			log.Fatalf("failed to run %T: %s", r, err)
		}
	}

	addr, ok := conf[runner.DiscoveryAddressKey]
	if !ok {
		log.Fatalf("unable to discover %s", runner.DiscoveryAddressKey)
	}

	// ctx, err := static.New("_fixtures/")
	// if err != nil {
	// 	log.Fatalf("failed to initialize context: %s", err)
	// }
	ctx, err := xds.New(addr)
	if err != nil {
		log.Fatalf("failed to initialize context: %s", err)
	}

	var results []types.AuditResult
	for _, auditor := range auditors {
		res, err := auditor.Audit(ctx)
		if err != nil {
			log.Printf("%s failed to run: %s", auditor.Name(), err)
		}
		results = append(results, res...)
	}

	for _, res := range results {
		fmt.Printf("%s: %s\n", res.Name, res.Description)
	}
}
