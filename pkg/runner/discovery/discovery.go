package discovery

import (
	"context"
	"fmt"

	"github.com/praetorian-inc/mithril/pkg/envoy"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/xds"
)

var Runner = runner.Runner{
	Name: "Discovery",
	Strategies: []runner.Strategy{
		&IstiodStrategy{},
		&IstioPilotStrategy{},
		&EnvoyConfigStrategy{},
	},
}

func verifyDiscoveryAddress(addr string) error {
	if addr == "" {
		return fmt.Errorf("empty addr")
	}
	cli, err := xds.NewClient(addr)
	if err != nil {
		return err
	}
	_, err = cli.Version(context.Background())
	return err
}

type IstiodStrategy struct{}

func (s *IstiodStrategy) Name() string {
	return "istiod"
}

func (s *IstiodStrategy) Run(input map[string]string) (map[string]string, error) {
	addr := fmt.Sprintf("istiod.%s.svc.cluster.local:15010", input[runner.IstioNamespaceKey])
	return map[string]string{
		runner.DiscoveryAddressKey: addr,
	}, nil
}

func (s *IstiodStrategy) Verify(input map[string]string) error {
	return verifyDiscoveryAddress(input[runner.DiscoveryAddressKey])
}

type IstioPilotStrategy struct{}

func (s *IstioPilotStrategy) Name() string {
	return "istio-pilot"
}

func (s *IstioPilotStrategy) Run(input map[string]string) (map[string]string, error) {
	addr := fmt.Sprintf("istio-pilot.%s.svc.cluster.local:15010", input[runner.IstioNamespaceKey])
	return map[string]string{
		runner.DiscoveryAddressKey: addr,
	}, nil
}

func (s *IstioPilotStrategy) Verify(input map[string]string) error {
	return verifyDiscoveryAddress(input[runner.DiscoveryAddressKey])
}

type EnvoyConfigStrategy struct{}

func (s *EnvoyConfigStrategy) Name() string {
	return "envoy"
}

func (s *EnvoyConfigStrategy) Run(map[string]string) (map[string]string, error) {
	// curl -s 127.0.0.1:15000/server_info | jq -r .node.metadata.PROXY_CONFIG.discoveryAddress
	ec, err := envoy.RetrieveConfig("http://localhost:15000/config_dump")
	if err != nil {
		return nil, err
	}
	address, err := ec.DiscoveryAddress()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		runner.DiscoveryAddressKey: address,
	}, nil
}

func (s *EnvoyConfigStrategy) Verify(input map[string]string) error {
	return verifyDiscoveryAddress(input[runner.DiscoveryAddressKey])
}
