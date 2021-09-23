package namespace

import (
	"fmt"

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
	// istiod.istio-system.svc.cluster.local:15010
	return map[string]string{
		runner.IstioNamespaceKey: "istio-system",
	}, nil
}

func (s *EnvoyStrategy) Verify(input map[string]string) error {
	return verifyIstioNamespace(input[runner.IstioNamespaceKey])
}
