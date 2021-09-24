package discovery

import (
	"fmt"
	"github.com/praetorian-inc/mithril/pkg/envoy"
	"log"

	"github.com/praetorian-inc/mithril/pkg/runner"
)

var Runner = runner.Runner{
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
	return nil
}

type DefaultStrategy struct{}

func (s *DefaultStrategy) Run(input map[string]string) (map[string]string, error) {
	log.Println("Attempting to use default discovery endpoint...")

	// istiod.istio-system.svc.cluster.local:15010
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
	log.Println("Attempting to retrieve discovery endpoint using Envoy admin panel...")

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
