// Copyright 2021 Praetorian Security, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package namespace implements a runner to locate the kubernetes namespace
// associated with the istio control plane. to accomplish this, it comes
// equipped with the following strategies:
//
// DefaultStrategy:
//    by default, the istio control plane is set up in a namespace called
//    `istio-system`. this strategy merely attempts to use this as the target
//    namespace.
//
// EnvoyStrategy:
//    given access to `127.0.0.1:15000/server_info`, the envoy debug service, a
//    client can extract the discovery address, which contains the istio namespace.
//
package namespace

import (
	"fmt"
	"regexp"

	"github.com/praetorian-inc/mithril/pkg/envoy"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/types"
)

// Runner defines the list of strategies to use to discover information about
// the Istio namespace.
var Runner = runner.Runner{
	Name: "Namespace",
	Strategies: []runner.Strategy{
		&envoyStrategy{},
		&defaultStrategy{},
	},
}

func verifyIstioNamespace(addr string) error {
	if addr == "" {
		return fmt.Errorf("empty addr")
	}
	return nil
}

type defaultStrategy struct{}

// Name returns the strategy name for reporting purposes.
func (s *defaultStrategy) Name() string {
	return "default"
}

// Run executes the default strategy and populates the Discovery type's
// IstioNamespace if it can verify the results.
func (s *defaultStrategy) Run(input *types.Discovery) error {
	// istiod.istio-system.svc.cluster.local:15010
	ns := "istio-system"

	err := verifyIstioNamespace(ns)
	if err != nil {
		return err
	}

	input.IstioNamespace = ns

	return nil
}

type envoyStrategy struct{}

// Name returns the strategy name for reporting purposes.
func (s *envoyStrategy) Name() string {
	return "envoy"
}

// Run executes the envoy strategy and populates the Discovery type's
// IstioNamespace if it can verify the results.
func (s *envoyStrategy) Run(input *types.Discovery) error {
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
