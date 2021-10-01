package install

import (
	log "github.com/sirupsen/logrus"

	"github.com/praetorian-inc/mithril/auditors"
	"github.com/praetorian-inc/mithril/pkg/types"
)

func init() {
	auditors.Register(&Auditor{})
}

type Auditor struct{}

func (a *Auditor) Name() string {
	return "Third Party Tokens"
}

func (a *Auditor) Audit(_ types.Discovery, resources types.Resources) ([]types.AuditResult, error) {
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

	if foundSidecar && policy != "third-party-jwt" {
		log.WithFields(log.Fields{
			"auditor": a.Name(),
			"policy":  policy,
		}).Info("found jwt policy")

		results = append(results, types.AuditResult{
			Name:        "Weak Service Account Authentication",
			Description: "JWT policy not set to third-party-jwt",
		})
	}

	return results, nil
}
