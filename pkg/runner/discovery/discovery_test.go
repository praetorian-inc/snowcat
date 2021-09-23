package discovery

import (
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/praetorian-inc/mithril/pkg/envoy"
	"testing"
)

const testJson = `
{
  "configs": [
    {
      "@type": "type.googleapis.com/envoy.admin.v3.BootstrapConfigDump",
      "bootstrap": {
        "node": {
          "cluster": "httpbin2.default",
          "id": "sidecar~10.4.2.16~httpbin2-5cb469675d-x2m7g.default~default.svc.cluster.local",
          "metadata": {
            "CLUSTER_ID": "Kubernetes",
            "INSTANCE_IPS": "10.4.2.16",
            "ISTIO_PROXY_SHA": "istio-proxy:4b528a87271e841bd64daf841a1a384ed4fcac68",
            "LABELS": {
              "app": "httpbin2",
              "istio.io/rev": "default",
              "pod-template-hash": "5cb469675d",
              "security.istio.io/tlsMode": "istio",
              "service.istio.io/canonical-name": "httpbin2",
              "service.istio.io/canonical-revision": "v1",
              "version": "v1"
            },
            "MESH_ID": "cluster.local",
            "NAME": "httpbin2-5cb469675d-x2m7g",
            "PLATFORM_METADATA": {
              "gcp_gce_instance_id": "7955581888538677628",
              "gcp_gke_cluster_name": "cluster-1",
              "gcp_gke_cluster_url": "https://container.googleapis.com/v1/projects/service-mithril/locations/us-central1-c/clusters/cluster-1",
              "gcp_location": "us-central1-c",
              "gcp_project": "service-mithril",
              "gcp_project_number": "199906641779"
            },
            "POD_PORTS": "[{\"containerPort\":80,\"protocol\":\"TCP\"}]",
            "PROXY_CONFIG": {
              "binaryPath": "/usr/local/bin/envoy",
              "concurrency": 2,
              "configPath": "./etc/istio/proxy",
              "controlPlaneAuthPolicy": "MUTUAL_TLS",
              "discoveryAddress": "istiod.istio-system.svc:15012",
              "drainDuration": "45s",
              "parentShutdownDuration": "60s",
              "proxyAdminPort": 15000,
              "serviceCluster": "httpbin2.default",
              "statNameLength": 189,
              "statusPort": 15020,
              "terminationDrainDuration": "5s",
              "tracing": {
                "zipkin": {
                  "address": "zipkin.istio-system:9411"
                }
              }
            },
            "SERVICE_ACCOUNT": "httpbin2",
            "kubectl.kubernetes.io/restartedAt": "2021-09-21T17:10:50Z"
          }
        }
      }
    }
  ]
}
`

func TestRetrieveEnvoyConfig(t *testing.T) {
	json := []byte(testJson)
	// Removing the URL check for now since you need an envoy instance hosted to test this
	//ec, _ := envoy.RetrieveConfig("http://localhost:15000/config_dump")
	ec, _ := envoy.LoadConfig(json)
	address, _ := ec.DiscoveryAddress()
	fmt.Println(address)
	assert.Equal(t, "istiod.istio-system.svc:15012", address)
}
