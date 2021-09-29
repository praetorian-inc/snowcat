package destinationrule

import (
	"fmt"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/types"
)

func init() {
	auditors.Register(&Auditor{})
}

type Auditor struct{}

func (a *Auditor) Name() string {
	return "TLS Validation in Destination Rule"
}

// Unsafe if TLS is SIMPLE and missing CA certificates
func isClientTLSSettingSafe(tls *networkingv1alpha3.ClientTLSSettings) bool {
	return tls == nil || tls.Mode.String() != "SIMPLE" || tls.CaCertificates != ""
}

func (a *Auditor) Audit(_ types.Discovery, resources types.Resources) ([]types.AuditResult, error) {
	var results []types.AuditResult

	for _, rule := range resources.DestinationRules {
		if !isClientTLSSettingSafe(rule.Spec.TrafficPolicy.Tls) {
			results = append(results, types.AuditResult{
				Name:        a.Name(),
				Description: fmt.Sprintf("%s rule missing CA certificates", rule.Name),
			})
		}
		for _, policy := range rule.Spec.TrafficPolicy.PortLevelSettings {
			if !isClientTLSSettingSafe(policy.Tls) {
				results = append(results, types.AuditResult{
					Name:        a.Name(),
					Description: fmt.Sprintf("%s rule missing CA certificates in traffic policy", rule.Name),
				})
			}
		}
	}
	return results, nil
}
