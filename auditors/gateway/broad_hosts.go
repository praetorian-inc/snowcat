// Package gateway provides auditor implementations that analyze
// Istio Gateways.
package gateway

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
	return "Overly Broad Gateway Hosts"
}

func (a *auditor) Audit(_ types.Discovery, resources types.Resources) ([]types.AuditResult, error) {
	var results []types.AuditResult

	for _, gateway := range resources.Gateways {
		for _, server := range gateway.Spec.Servers {
			for _, host := range server.Hosts {
				if host == "*" {
					results = append(results, types.AuditResult{
						Name:        a.Name(),
						Resource:    gateway.Namespace + ":" + gateway.Name,
						Description: fmt.Sprintf("%s host gateway is too broad (wildcard host allowed)", gateway.Spec.Selector["istio"]),
					})
				}
			}
		}
	}

	return results, nil
}
