package gateway

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
	return "Overly Broad Gateway Hosts"
}

func (a *Auditor) Audit(c types.IstioContext) ([]types.AuditResult, error) {
	var results []types.AuditResult

	gateways, err := c.Gateways()
	if err != nil {
		return nil, fmt.Errorf("failed to get gateways: %w", err)
	}
	for _, gateway := range gateways {
		for _, server := range gateway.Spec.Servers {
			for _, host := range server.Hosts {
				if host == "*" {
					results = append(results, types.AuditResult{
						Name:        a.Name(),
						Description: fmt.Sprintf("%s host gateway is too broad (wildcard host allowed)", gateway),
					})
				}
			}
		}
	}

	return results, nil
}
