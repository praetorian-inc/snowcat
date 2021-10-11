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

// Package peerauth provides auditor implementations that analyze
// Istio PeerAuthentication policies.
package peerauth

import (
	"fmt"

	"github.com/praetorian-inc/snowcat/auditors"
	"github.com/praetorian-inc/snowcat/pkg/types"
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
