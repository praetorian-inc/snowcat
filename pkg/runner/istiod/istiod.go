package istiod

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"

	"github.com/praetorian-inc/mithril/pkg/debugz"
	"github.com/praetorian-inc/mithril/pkg/kubelet"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/types"
	"github.com/praetorian-inc/mithril/pkg/xds"
)

var Runner = runner.Runner{
	Name: "Istio Control Plane",
	Strategies: []runner.Strategy{
		&KubeletStrategy{},
		// &portScanStrategy{},
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
	if err != nil {
		return false
	}
	return true
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

type KubeletStrategy struct{}

func (s *KubeletStrategy) Name() string {
	return "kubelet"
}

func (s *KubeletStrategy) Run(input *types.Discovery) error {
	ctx := context.Background()

	var ips []string

	for _, addr := range input.KubeletAddresses {
		log.Printf("fetching pods from %s kubelet", addr)
		k, err := kubelet.NewClient(addr)
		if err != nil {
			log.Printf("failed to connect to %s kubelet: %s", addr, err)
			continue
		}

		pods, err := k.Pods(ctx)
		if err != nil {
			log.Printf("failed to list pods from %s kubelet: %s", addr, err)
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

	input.IstiodIPs = ips
	return nil
}
