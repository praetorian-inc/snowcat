package auth

import (
	"github.com/praetorian-inc/mithril/pkg/types"
)

type Auditor struct {
}

func (a *Auditor) Audit(c types.IstioContext) ([]types.AuditResult, error) {
	return []types.AuditResult{}, nil
}
