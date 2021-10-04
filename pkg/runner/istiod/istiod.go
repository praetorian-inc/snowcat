/* package istiod implements a runner to locate services associated with the
   istio control plane components. to accomplish this, it comes equipped with
   the following strategies:

KubeletStrategy:
   if provided with an ip:port combination known to be running the kubelet api,
   this strategy can query for running pods, check them for istio related
   labels, and determine whether or not they are running the debug/discovery
   services

IstiodStrategy:
   if provided with the istio namespace, it will attempt to locate the
   debug/discovery service at `istiod.{namespace}.svc.cluster.local`

IstioPilotStrategy:
   if provided with the istio namespace, it will attempt to locate the
   debug/discovery services at `istio-pilot.{namespace}.svc.cluster.local`

EnvoyConfigStrategy:
   the envoy configuration strategy will attempt to connect to
   `http://localhost:15000/config_dump` to extract the location of the discovery
   address
*/
package istiod

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"

	"github.com/praetorian-inc/mithril/pkg/debugz"
	"github.com/praetorian-inc/mithril/pkg/envoy"
	"github.com/praetorian-inc/mithril/pkg/kubelet"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/types"
	"github.com/praetorian-inc/mithril/pkg/xds"
)

/* Runner provides the list of strategies to use to gather information about the
   istio control plane components
*/
var Runner = runner.Runner{
	Name: "Istio Control Plane",
	Strategies: []runner.Strategy{
		&kubeletStrategy{},
		&istiodStrategy{},
		&istioPilotStrategy{},
		&envoyConfigStrategy{},
	},
}

func isRunning(pod v1.Pod) bool {
	return pod.Status.Phase == v1.PodRunning
}

func isIstiod(pod v1.Pod) bool {
	return pod.Labels["app"] == "istiod" ||
		pod.Labels["istio"] == "pilot" ||
		pod.Labels["operator.istio.io/component"] == "Pilot"
}

func hasDiscoveryService(host string) bool {
	c, err := xds.NewClient(host + ":15010")
	if err != nil {
		return false
	}
	c.Close()
	return true
}

func hasDebugService(host string) bool {
	_, err := debugz.NewClient(host + ":8080")
	return err == nil
}

func isDebugIstiod(host string) bool {
	if host == "" {
		return false
	}
	if hasDiscoveryService(host) {
		return true
	}
	if hasDebugService(host) {
		return true
	}
	return false
}

type kubeletStrategy struct{}

/* Name returns the strategy name for reporting purposes.

   required by the Strategy interface
*/
func (s *kubeletStrategy) Name() string {
	return "kubelet"
}

/* Run executes the kubelet strategy and populates the Discovery type's
   DiscoveryAddress and DebugzAddress if it can verify the results.

   required by the Strategy interface
*/
func (s *kubeletStrategy) Run(input *types.Discovery) error {
	ctx := context.Background()

	var ips []string

	for _, addr := range input.KubeletAddresses {
		k, err := kubelet.NewClient(addr)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": addr,
				"err":  err,
			}).Warn("failed to connect to kubelet")
			continue
		}

		pods, err := k.Pods(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": addr,
				"err":  err,
			}).Warn("failed to list pods from kubelet")
			continue
		}

		for _, pod := range pods {
			ip := pod.Status.PodIP
			if isRunning(pod) && isIstiod(pod) && isDebugIstiod(ip) {
				ips = append(ips, ip)
			}
		}
	}

	if len(ips) == 0 {
		return fmt.Errorf("failed to find istiod")
	}

	var found bool
	for _, ip := range ips {
		if hasDiscoveryService(ip) {
			input.DiscoveryAddress = ip + ":15010"
			found = true
		}
		if hasDebugService(ip) {
			input.DebugzAddress = ip + ":8080"
			found = true
		}
		if found {
			break
		}
	}
	return nil
}

type istiodStrategy struct{}

/* Name returns the strategy name for reporting purposes.

   required by the Strategy interface
*/
func (s *istiodStrategy) Name() string {
	return "istiod"
}

/* Run executes the istiod strategy and populates the Discovery type's
   DiscoveryAddress and DebugzAddress if it can verify the results.

   required by the Strategy interface
*/
func (s *istiodStrategy) Run(input *types.Discovery) error {
	if input.IstioNamespace == "" {
		return fmt.Errorf("istio namespace required")
	}
	addr := fmt.Sprintf("istiod.%s.svc.cluster.local", input.IstioNamespace)
	if hasDiscoveryService(addr) {
		input.DiscoveryAddress = addr + ":15010"
	}
	if hasDebugService(addr) {
		input.DebugzAddress = addr + ":8080"
	}
	return nil
}

type istioPilotStrategy struct{}

/* Name returns the strategy name for reporting purposes.

   required by the Strategy interface
*/
func (s *istioPilotStrategy) Name() string {
	return "istio-pilot"
}

/* Run executes the istio-pilot strategy and populates the Discovery type's
   DiscoveryAddress and DebugzAddress if it can verify the results.

   required by the Strategy interface
*/
func (s *istioPilotStrategy) Run(input *types.Discovery) error {
	if input.IstioNamespace == "" {
		return fmt.Errorf("istio namespace required")
	}
	addr := fmt.Sprintf("istio-pilot.%s.svc.cluster.local", input.IstioNamespace)
	if hasDiscoveryService(addr) {
		input.DiscoveryAddress = addr + ":15010"
	}
	if hasDebugService(addr) {
		input.DebugzAddress = addr + ":8080"
	}
	return nil
}

type envoyConfigStrategy struct{}

/* Name returns the strategy name for reporting purposes.

   required by the Strategy interface
*/
func (s *envoyConfigStrategy) Name() string {
	return "envoy"
}

/* Run executes the envoy config strategy and populates the Discovery type's
   DiscoveryAddress if it can verify the results.

   required by the Strategy interface
*/
func (s *envoyConfigStrategy) Run(input *types.Discovery) error {
	ec, err := envoy.RetrieveConfig("http://localhost:15000/config_dump")
	if err != nil {
		return err
	}
	addr, err := ec.DiscoveryAddress()
	if err != nil {
		return err
	}

	c, err := xds.NewClient(addr)
	if err != nil {
		return err
	}
	c.Close()

	input.DiscoveryAddress = addr
	return nil
}
