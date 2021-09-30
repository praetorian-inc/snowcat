package auth

import (
	"fmt"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/types"
)

func init() {
	auditors.Register(&Auditor{})
}

type Auditor struct{}

func (a *Auditor) Name() string {
	return "Mututal TLS"
}

func (a *Auditor) Audit(_ types.Discovery, resources types.Resources) ([]types.AuditResult, error) {
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
