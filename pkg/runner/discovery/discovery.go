package discovery

import (
	"fmt"

	"github.com/praetorian-inc/mithril/pkg/runner"
)

var Runner = runner.Runner{
	Strategies: []runner.Strategy{
		&DefaultStrategy{},
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
	// istiod.istio-system.svc.cluster.local:15010
	addr := fmt.Sprintf("istiod.%s.svc.cluster.local:15010", input[runner.IstioNamespaceKey])
	return map[string]string{
		runner.DiscoveryAddressKey: addr,
	}, nil
}

func (s *DefaultStrategy) Verify(input map[string]string) error {
	return verifyDiscoveryAddress(input[runner.DiscoveryAddressKey])
}

type EnvoyConfigStrategy struct{}

func (s *EnvoyConfigStrategy) Run(map[string]string) (map[string]string, error) {
	// curl -s 127.0.0.1:15000/server_info | jq -r .node.metadata.PROXY_CONFIG.discoveryAddress
	return map[string]string{}, nil
}

func (s *EnvoyConfigStrategy) Verify(input map[string]string) error {
	return verifyDiscoveryAddress(input[runner.DiscoveryAddressKey])
}
