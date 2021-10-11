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

// Package destinationrule provides auditor implementations that analyze
// Istio DestinationRules.
package destinationrule

import (
	"fmt"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"

	"github.com/praetorian-inc/snowcat/auditors"
	"github.com/praetorian-inc/snowcat/pkg/types"
)

func init() {
	auditors.Register(&auditor{})
}

type auditor struct{}

func (a *auditor) Name() string {
	return "TLS Validation in Destination Rule"
}

// Unsafe if TLS is SIMPLE and missing CA certificates
func isClientTLSSettingSafe(tls *networkingv1alpha3.ClientTLSSettings) bool {
	return tls == nil || tls.Mode.String() != "SIMPLE" || tls.CaCertificates != ""
}

func (a *auditor) Audit(_ types.Discovery, resources types.Resources) ([]types.AuditResult, error) {
	var results []types.AuditResult

	for _, rule := range resources.DestinationRules {
		if !isClientTLSSettingSafe(rule.Spec.TrafficPolicy.Tls) {
			results = append(results, types.AuditResult{
				Name:        a.Name(),
				Resource:    rule.Namespace + ":" + rule.Name,
				Description: fmt.Sprintf("%s rule missing CA certificates", rule.Name),
			})
		}
		for _, policy := range rule.Spec.TrafficPolicy.PortLevelSettings {
			if !isClientTLSSettingSafe(policy.Tls) {
				results = append(results, types.AuditResult{
					Name:        a.Name(),
					Resource:    rule.Namespace + ":" + rule.Name,
					Description: fmt.Sprintf("%s rule missing CA certificates in traffic policy", rule.Name),
				})
			}
		}
	}
	return results, nil
}
