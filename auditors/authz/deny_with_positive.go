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

package authz

import (
	"fmt"

	apiv1beta "istio.io/api/security/v1beta1"
	security "istio.io/client-go/pkg/apis/security/v1beta1"

	"github.com/praetorian-inc/mesh-hunter/auditors"
	"github.com/praetorian-inc/mesh-hunter/pkg/types"
)

func init() {
	auditors.Register(&denyWithPositiveAuditor{})
}

type denyWithPositiveAuditor struct{}

func (a *denyWithPositiveAuditor) Name() string {
	return "Deny with Positive Match"
}

func (a *denyWithPositiveAuditor) Audit(_ types.Discovery, resources types.Resources) ([]types.AuditResult, error) {
	var results []types.AuditResult

	for _, policy := range resources.AuthorizationPolicies {
		if policy.Spec.Action != apiv1beta.AuthorizationPolicy_DENY {
			continue
		}

		offendingRules := evalDenyPolicy(policy)

		if len(offendingRules) > 0 {
			results = append(results, types.AuditResult{
				Name:        a.Name(),
				Resource:    policy.Namespace + ":" + policy.Name,
				Description: fmt.Sprintf("discovered deny policy with positive matchers in %s", policy.Name),
			})
		}
	}

	return results, nil
}

func evalDenyPolicy(policy security.AuthorizationPolicy) []apiv1beta.Rule {
	var offendingRules []apiv1beta.Rule

	rules := policy.Spec.Rules

	for _, rule := range rules {
		fromRules := rule.From

		for _, f := range fromRules {
			source := f.Source

			if source.IpBlocks != nil ||
				source.Namespaces != nil ||
				source.Principals != nil ||
				source.RemoteIpBlocks != nil ||
				source.RequestPrincipals != nil {
				offendingRules = append(offendingRules, *rule)
			}
		}

		toRules := rule.To

		for _, t := range toRules {
			operation := t.Operation

			if operation.Hosts != nil ||
				operation.Methods != nil ||
				operation.Paths != nil ||
				operation.Ports != nil {
				offendingRules = append(offendingRules, *rule)
			}
		}
	}

	return offendingRules
}
