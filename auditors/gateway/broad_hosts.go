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

// Package gateway provides auditor implementations that analyze
// Istio Gateways.
package gateway

import (
	"fmt"

	"github.com/praetorian-inc/mesh-hunter/auditors"
	"github.com/praetorian-inc/mesh-hunter/pkg/types"
)

func init() {
	auditors.Register(&auditor{})
}

type auditor struct{}

func (a *auditor) Name() string {
	return "Overly Broad Gateway Hosts"
}

func (a *auditor) Audit(_ types.Discovery, resources types.Resources) ([]types.AuditResult, error) {
	var results []types.AuditResult

	for _, gateway := range resources.Gateways {
		for _, server := range gateway.Spec.Servers {
			for _, host := range server.Hosts {
				if host == "*" {
					results = append(results, types.AuditResult{
						Name:        a.Name(),
						Resource:    gateway.Namespace + ":" + gateway.Name,
						Description: fmt.Sprintf("%s host gateway is too broad (wildcard host allowed)", gateway.Spec.Selector["istio"]),
					})
				}
			}
		}
	}

	return results, nil
}
