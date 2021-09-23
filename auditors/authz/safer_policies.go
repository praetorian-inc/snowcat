package authz

import (
	"fmt"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/types"

	security "istio.io/client-go/pkg/apis/security/v1beta1"

	apiv1beta "pkg.go.dev/istio.io/api/security/v1beta1"
)

const (
	ALLOW = 0
	DENY  = 1
)

func init() {
	auditors.Register(&Auditor{})
}

type Auditor struct{}

func (a *Auditor) Name() string {
	return "Safer Authorization Policies"
}

func (a *Auditor) Audit(c types.IstioContext) ([]types.AuditResult, error) {
	var results []types.AuditResult

	policies, err := c.AuthorizationPolicies()
	if err != nil {
		return nil, fmt.Errorf("failed to get authz policies: %w", err)
	}

	for _, policy := range policies {
		switch action := policy.Spec.Action; action {
		case ALLOW:
			offendingRules := evalAllowPolicy(policy)

			if offendingRules != nil {
				results = append(results, types.AuditResult{
					Name:        a.Name(),
					Description: fmt.Sprintf("discovered allow policy with negative matchers in %s", policy.Name),
				})
			}
		case DENY:
			offendingRules := evalDenyPolicy(policy)

			if offendingRules != nil {
				results = append(results, types.AuditResult{
					Name:        a.Name(),
					Description: fmt.Sprintf("discovered deny policy with positive matchers in %s", policy.Name),
				})
			}
		default:
			continue
		}
	}

	return results, nil
}

func evalAllowPolicy(policy security.AuthorizationPolicy) []apiv1beta.Rule {
	var offendingRules []apiv1beta.Rule

	rules := policy.Spec.Rules

	for _, rule := range rules {
		fromRules := rule.From

		for _, f := range fromRules {
			source := f.Source

			if source.NotIpBlocks != nil ||
				source.NotNamespaces != nil ||
				source.NotPrincipals != nil ||
				source.NotRemoteIpBlocks != nil ||
				source.NotRequestPrincipals != nil {

				offendingRules = append(offendingRules, rule)
				continue // this shadows errors in the to rules, we may not want this
			}
		}

		toRules := rule.To

		for _, t := range toRules {
			operation := t.Operation

			if operation.NotHosts != nil ||
				operation.NotMethods != nil ||
				operation.NotPaths != nil ||
				operation.NotPorts != nil {

				offendingRules = append(offendingRules, rule)
				continue
			}
		}
	}

	return offendingRules
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

				offendingRules = append(offendingRules, rule)
				continue // this shadows errors in the to rules, we may not want this
			}
		}

		toRules := rule.To

		for _, t := range toRules {
			operation := t.Operation

			if operation.Hosts != nil ||
				operation.Methods != nil ||
				operation.Paths != nil ||
				operation.Ports != nil {

				offendingRules = append(offendingRules, rule)
				continue
			}
		}
	}

	return offendingRules
}
