package authz

import (
	"fmt"

	apiv1beta "istio.io/api/security/v1beta1"
	security "istio.io/client-go/pkg/apis/security/v1beta1"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/types"
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
