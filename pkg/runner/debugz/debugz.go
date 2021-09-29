package debugz

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"

	"github.com/praetorian-inc/mithril/pkg/debugz"
	"github.com/praetorian-inc/mithril/pkg/kubelet"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/xds"
)

var Runner = runner.Runner{
	Name: "Istio Control Plane",
	Strategies: []runner.Strategy{
		&KubeletStrategy{},
		// &portScanStrategy{},
	},
}

func isIstiod(pod v1.Pod) bool {
	return pod.Labels["app"] == "istiod" ||
		pod.Labels["istio"] == "pilot" ||
		pod.Labels["operator.istio.io/component"] == "Pilot"
}

func isDebugIstiod(host string) bool {
	if host == "" {
		return false
	}

	log.Printf("checking %s for debug ports", host)

	c, err := xds.NewClient(host + ":15010")
	if err != nil {
		return false
	}
	c.Close()

	_, err = debugz.NewClient(host + ":8080")
	if err != nil {
		return false
	}

	return true
}

type KubeletStrategy struct{}

func (s *KubeletStrategy) Name() string {
	return "kubelet"
}

func (s *KubeletStrategy) Run(input map[string]string) (map[string]string, error) {
	ctx := context.Background()

	var istios []string

	// TODO: fetch this from runner input
	kubelets := []string{"10.48.0.1:10255", "10.48.1.1:10255", "10.48.2.1:10255", "10.48.3.1:10255"}

	for _, addr := range kubelets {
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
			if (pod.Status.Phase == v1.PodRunning) && isIstiod(pod) && isDebugIstiod(ip) {
				istios = append(istios, ip)
			}
		}
	}

	if len(istios) == 0 {
		return nil, fmt.Errorf("failed to find istiod")
	}

	for _, host := range istios {
		log.Printf("%s", host)
	}

	// TODO: set struct field
	/*
		return map[string]string{
			runner.IstioIPsKey: "istios",
		}, nil
	*/
	return map[string]string{}, nil
}

func (s *KubeletStrategy) Verify(input map[string]string) error {
	return nil
}
