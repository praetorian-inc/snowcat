// Package version provides auditor implementations that analyze
// the version of Istio for known CVEs.
package version

import (
	"fmt"

	"github.com/praetorian-inc/mithril/pkg/knownvulns"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/types"
)

func init() {
	auditors.Register(&auditor{})
}

type auditor struct{}

func (a *auditor) Name() string {
	return "Known Vulnerable Version"
}

func (a *auditor) Audit(disco types.Discovery, _ types.Resources) ([]types.AuditResult, error) {
	var results []types.AuditResult

	if disco.IstioVersion == "" {
		return nil, fmt.Errorf("version required")
	}

	vulns, err := knownvulns.GetVulnsForVersion(disco.IstioVersion)
	if err != nil {
		return nil, fmt.Errorf("error retrieving vulns for version: %w", err)
	}

	for _, vuln := range vulns {
		results = append(results, types.AuditResult{
			Name:        a.Name(),
			Resource:    "Version " + disco.IstioVersion,
			Description: fmt.Sprintf("Vulnerable to %s (Impact Score %s) - more details at %s", vuln.DisclosureID, vuln.ImpactScore, vuln.DisclosureURL),
		})
	}

	return results, nil
}
