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

// Package version provides auditor implementations that analyze
// the version of Istio for known CVEs.
package version

import (
	"fmt"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/knownvulns"
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
			Name:     a.Name(),
			Resource: "Version " + disco.IstioVersion,
			Description: fmt.Sprintf("Vulnerable to %s (Impact Score %s) - more details at %s",
				vuln.DisclosureID, vuln.ImpactScore, vuln.DisclosureURL),
		})
	}

	return results, nil
}
