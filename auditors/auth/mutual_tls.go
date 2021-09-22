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

func (a *Auditor) Audit(c types.IstioContext) ([]types.AuditResult, error) {
	var results []types.AuditResult

	namespaces, err := c.Namespaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}
	namespaceSafety := make(map[string]bool)
	for _, ns := range namespaces {
		namespaceSafety[ns] = false
	}

	policies, err := c.PeerAuthentications()
	if err != nil {
		return nil, fmt.Errorf("failed to get peer authentication policies: %w", err)
	}
	for _, policy := range policies {
		if policy.Spec.Mtls.Mode.String() == "STRICT" {
			namespaceSafety[policy.Namespace] = true
		}
	}

	for ns, safe := range namespaceSafety {
		if !safe {
			results = append(results, types.AuditResult{
				Name:        a.Name(),
				Description: fmt.Sprintf("%s namespace missing PeerAuthentication policy", ns),
			})
		}
	}
	return results, nil
}
