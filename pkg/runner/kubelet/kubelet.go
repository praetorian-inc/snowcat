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

// Package kubelet implements a runner to locate services associated with the
// current cluster's kubelet api. to accomplish this, it comes equipped with
// the following strategies:
//
// DefaultGatewayStrategy:
//    this strategy will attempt to locate the current pod's default gateway. from
//    there, it will scan every subnet for the same address to look for http
//    services associated with the kubelet api.
//
//    i.e. default gateway = 192.168.2.1 will produce a scan for 192.168.{0-255}.1:10255
package kubelet

import (
	"net"
	"time"

	"github.com/jackpal/gateway"
	log "github.com/sirupsen/logrus"

	"github.com/praetorian-inc/mesh-hunter/pkg/kubelet"
	"github.com/praetorian-inc/mesh-hunter/pkg/netscan"
	"github.com/praetorian-inc/mesh-hunter/pkg/runner"
	"github.com/praetorian-inc/mesh-hunter/pkg/types"
)

// Runner defines the list of strategies to use to discover information about
// the Kubelet read-only API on nodes in the cluster.
var Runner = runner.Runner{
	Name: "Kubelet",
	Strategies: []runner.Strategy{
		&defaultGatewayStrategy{},
	},
}

func verifyKubeletAPI(addr string) bool {
	_, err := kubelet.NewClient(addr)
	return err == nil
}

type defaultGatewayStrategy struct{}

// Name returns the strategy name for reporting purposes.
func (s *defaultGatewayStrategy) Name() string {
	return "default-gateway"
}

// Run executes the default gateway strategy and populates the Discovery type's
// KubeletAddress if it can verify the results.
func (s *defaultGatewayStrategy) Run(input *types.Discovery) error {
	var hosts []string

	log.Info("attempting to locate default gateway")

	// this could be an issue depending on the platform the pod is using
	gateway, err := gateway.DiscoverGateway()
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"gateway": gateway.String(),
	}).Info("discovered default gateway")

	if ip4 := gateway.To4(); ip4 != nil {
		for i := 0; i < 256; i++ {
			ip := net.IPv4(ip4[0], ip4[1], byte(i), ip4[3])
			hosts = append(hosts, ip.String())
		}
	}

	log.Info("scanning for additional potential gateways using HTTP scanner")

	scanner, err := netscan.New(netscan.ModeHTTP, hosts, []string{"10255"})
	if err != nil {
		return err
	}

	var results []string

	for addr := range scanner.Scan(500 * time.Millisecond) {
		if verifyKubeletAPI(addr) {
			log.WithFields(log.Fields{
				"addr": addr,
			}).Debug("discovered kubelet api")

			results = append(results, addr)
		}
	}

	log.WithFields(log.Fields{
		"kubeletAddresses": results,
	}).Debug("resulting kubelet apis")

	input.KubeletAddresses = results
	return nil
}
