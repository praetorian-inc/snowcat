package namespace

import (
	"fmt"
	"regexp"

	"github.com/praetorian-inc/mithril/pkg/envoy"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/types"
)

var Runner = runner.Runner{
	Name: "Namespace",
	Strategies: []runner.Strategy{
		&EnvoyStrategy{},
		&DefaultStrategy{},
	},
}

func verifyIstioNamespace(addr string) error {
	if addr == "" {
		return fmt.Errorf("empty addr")
	}
	return nil
}

type DefaultStrategy struct{}

func (s *DefaultStrategy) Name() string {
	return "default"
}

func (s *DefaultStrategy) Run(input *types.Discovery) error {
	// istiod.istio-system.svc.cluster.local:15010
	ns := "istio-system"

	err := verifyIstioNamespace(ns)
	if err != nil {
		return err
	}

	input.IstioNamespace = ns

	return nil
}

type EnvoyStrategy struct{}

func (s *EnvoyStrategy) Name() string {
	return "envoy"
}

func (s *EnvoyStrategy) Run(input *types.Discovery) error {
	// curl -s 127.0.0.1:15000/server_info | jq -r .node.metadata.PROXY_CONFIG.discoveryAddress
	ec, err := envoy.RetrieveConfig("http://localhost:15000/config_dump")
	if err != nil {
		return err
	}

	address, err := ec.DiscoveryAddress()
	if err != nil {
		return err
	}

	r := regexp.MustCompile(`istiod\.(.*)\.svc.*:.*`)

	match := r.FindStringSubmatch(address)
	if match == nil {
		return fmt.Errorf("%s did not match expected regex of `istiod.<namespace>.svc:<port>`", address)
	}

	err = verifyIstioNamespace(match[1])
	if err != nil {
		return err
	}

	input.IstioNamespace = match[1]

	return nil
}
