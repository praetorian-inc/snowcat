package version

import (
	"fmt"
	"github.com/praetorian-inc/mithril/pkg/knownvulns"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/types"
)

func init() {
	auditors.Register(&Auditor{})
}

type Auditor struct{}

func (a *Auditor) Name() string {
	return "Known Vulnerable Version"
}

func (a *Auditor) Audit(c types.IstioContext) ([]types.AuditResult, error) {
	var results []types.AuditResult

	version, err := c.Version()
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	vulns, err := knownvulns.GetVulnsForVersion(version)
	if err != nil {
		return nil, fmt.Errorf("error retrieving vulns for version: %w", err)
	}

	for _, vuln := range vulns {
		results = append(results, types.AuditResult{
			Name:        a.Name(),
			Description: fmt.Sprintf("Vulnerable to %s (Impact Score %s) - more details at %s", vuln.DisclosureId, vuln.ImpactScore, vuln.DisclosureUrl),
		})
	}

	return results, nil
}
