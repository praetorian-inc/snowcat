package discovery

import (
	"fmt"

	"github.com/praetorian-inc/mithril/pkg/runner"
)

var Runner = runner.Runner{
	Strategies: []runner.Strategy{
		&StaticStrategy{},
		&DynamicStrategy{},
	},
}

func VerifyDiscoveryAddress(addr string) error {
	return nil
}

type StaticStrategy struct {
}

func (s *StaticStrategy) Run(input map[string]string) (map[string]string, error) {
	addr := fmt.Sprintf("istiod.%s.svc.cluster.local:15010", input[runner.IstioNamespaceKey])
	return map[string]string{
		runner.DiscoveryAddressKey: addr,
	}, VerifyDiscoveryAddress(addr)
}

type DynamicStrategy struct {
}

func (s *DynamicStrategy) Run(map[string]string) (map[string]string, error) {
	return map[string]string{}, nil
}
