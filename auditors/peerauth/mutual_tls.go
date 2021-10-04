// Package peerauth provides auditor implementations that analyze
// Istio PeerAuthentication policies.
package peerauth

import (
	"fmt"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/types"
)

func init() {
	auditors.Register(&auditor{})
}

type auditor struct{}

func (a *auditor) Name() string {
	return "Permissive Mututal TLS"
}

func (a *auditor) Audit(_ types.Discovery, resources types.Resources) ([]types.AuditResult, error) {
	var results []types.AuditResult

	namespaceSafety := make(map[string]bool)
	for _, ns := range resources.Namespaces {
		namespaceSafety[ns.Name] = false
	}

	for _, policy := range resources.PeerAuthentications {
		if policy.Spec.Mtls.Mode.String() == "STRICT" {
			namespaceSafety[policy.Namespace] = true
		}
	}

	for ns, safe := range namespaceSafety {
		if !safe {
			results = append(results, types.AuditResult{
				Name:        a.Name(),
				Resource:    ns,
				Description: fmt.Sprintf("%s namespace missing PeerAuthentication policy", ns),
			})
		}
	}
	return results, nil
}
