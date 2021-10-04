// Package install provides auditor implementations that analyze
// the IstioOperator and general control plane configurations.
package install

import (
	log "github.com/sirupsen/logrus"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/types"
)

func init() {
	auditors.Register(&auditor{})
}

type auditor struct{}

func (a *auditor) Name() string {
	return "Weak Service Account Authentication"
}

func (a *auditor) Audit(_ types.Discovery, resources types.Resources) ([]types.AuditResult, error) {
	var results []types.AuditResult

	var policy string
	var foundSidecar bool

	// Iterate over all pods with the sidecar.istio.io/status annotation
	for _, pod := range resources.Pods {
		if pod.Annotations["sidecar.istio.io/status"] == "" {
			continue
		}

		// Iterate over all containers, searching for the istio-proxy container
		for _, container := range pod.Spec.Containers {
			if container.Name != "istio-proxy" {
				continue
			}

			// Iterate over all istio-proxy env vars, searching for the JWT_POLICY
			for _, env := range container.Env {
				if env.Name == "JWT_POLICY" {
					policy = env.Value
				}
			}

			foundSidecar = true
		}
	}

	if !foundSidecar {
		log.Warn("no istio sidecars found")
	}

	if policy == "" {
		log.WithFields(log.Fields{
			"auditor": a.Name(),
		}).Info("found no active jwt policy")
	} else {
		log.WithFields(log.Fields{
			"auditor": a.Name(),
			"policy":  policy,
		}).Info("found jwt policy")
	}

	if foundSidecar && policy != "third-party-jwt" {
		results = append(results, types.AuditResult{
			Name:        a.Name(),
			Description: "JWT policy not set to third-party-jwt",
		})
	}

	return results, nil
}
