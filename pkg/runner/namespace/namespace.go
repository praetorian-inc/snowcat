package namespace

import (
	"fmt"
	"regexp"

	"github.com/praetorian-inc/mithril/pkg/envoy"
	"github.com/praetorian-inc/mithril/pkg/runner"
)

var Runner = runner.Runner{
	Strategies: []runner.Strategy{
		&DefaultStrategy{},
		&EnvoyStrategy{},
	},
}

func verifyIstioNamespace(addr string) error {
	if addr == "" {
		return fmt.Errorf("empty addr")
	}
	return nil
}

type DefaultStrategy struct{}

func (s *DefaultStrategy) Run(input map[string]string) (map[string]string, error) {
	// istiod.istio-system.svc.cluster.local:15010
	return map[string]string{
		runner.IstioNamespaceKey: "istio-system",
	}, nil
}

func (s *DefaultStrategy) Verify(input map[string]string) error {
	return verifyIstioNamespace(input[runner.IstioNamespaceKey])
}

type EnvoyStrategy struct{}

func (s *EnvoyStrategy) Run(input map[string]string) (map[string]string, error) {
	// curl -s 127.0.0.1:15000/server_info | jq -r .node.metadata.PROXY_CONFIG.discoveryAddress
	ec, err := envoy.RetrieveConfig("http://localhost:15000/config_dump")
	if err != nil {
		return nil, err
	}

	address, err := ec.DiscoveryAddress()
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(`istiod\.(.*)\.svc.*:.*`)

	match := r.FindStringSubmatch(address)
	if match == nil {
		return nil, fmt.Errorf("%s did not match expected regex of `istiod.<namespace>.svc:<port>`", address)
	}

	return map[string]string{
		runner.IstioNamespaceKey: match[1],
	}, nil
}

func (s *EnvoyStrategy) Verify(input map[string]string) error {
	return verifyIstioNamespace(input[runner.IstioNamespaceKey])
}
