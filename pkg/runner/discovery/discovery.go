package discovery

import (
	"context"
	"fmt"

	"github.com/praetorian-inc/mithril/pkg/envoy"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/types"
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

func (s *IstiodStrategy) Run(input *types.Discovery) error {
	if input.IstioNamespace == "" {
		return fmt.Errorf("istio namespace required")
	}
	addr := fmt.Sprintf("istiod.%s.svc.cluster.local:15010", input.IstioNamespace)
	if err := verifyDiscoveryAddress(addr); err != nil {
		return err
	}
	input.DiscoveryAddress = addr
	return nil
}

type IstioPilotStrategy struct{}

func (s *IstioPilotStrategy) Name() string {
	return "istio-pilot"
}

func (s *IstioPilotStrategy) Run(input *types.Discovery) error {
	if input.IstioNamespace == "" {
		return fmt.Errorf("istio namespace required")
	}
	addr := fmt.Sprintf("istio-pilot.%s.svc.cluster.local:15010", input.IstioNamespace)
	if err := verifyDiscoveryAddress(addr); err != nil {
		return err
	}
	input.DiscoveryAddress = addr
	return nil
}

type EnvoyConfigStrategy struct{}

func (s *EnvoyConfigStrategy) Name() string {
	return "envoy"
}

func (s *EnvoyConfigStrategy) Run(input *types.Discovery) error {
	ec, err := envoy.RetrieveConfig("http://localhost:15000/config_dump")
	if err != nil {
		return err
	}
	addr, err := ec.DiscoveryAddress()
	if err != nil {
		return err
	}
	if err := verifyDiscoveryAddress(addr); err != nil {
		return err
	}
	input.DiscoveryAddress = addr
	return nil
}
