package kubelet

import (
	"net"
	"time"

	"github.com/jackpal/gateway"

	"github.com/praetorian-inc/mithril/pkg/kubelet"
	"github.com/praetorian-inc/mithril/pkg/portscan"
	"github.com/praetorian-inc/mithril/pkg/runner"
	"github.com/praetorian-inc/mithril/pkg/types"
)

var Runner = runner.Runner{
	Name: "Kubelet",
	Strategies: []runner.Strategy{
		&DefaultGatewayStrategy{},
	},
}

func verifyKubeletAPI(addr string) bool {
	_, err := kubelet.NewClient(addr)
	return err == nil
}

type DefaultGatewayStrategy struct{}

func (s *DefaultGatewayStrategy) Name() string {
	return "default-gateway"
}

func (s *DefaultGatewayStrategy) Run(input *types.Discovery) error {
	var ips []net.IP

	gateway, err := gateway.DiscoverGateway()
	if err != nil {
		return err
	}

	if ip4 := gateway.To4(); ip4 != nil {
		for i := 1; i < 256; i++ {
			ips = append(ips, net.IPv4(ip4[0], ip4[1], byte(i), ip4[3]))
		}
	}

	scanner := portscan.New(ips, []int{10255})

	var results []string

	for addr := range scanner.Start(500 * time.Millisecond) {
		if verifyKubeletAPI(addr) {
			results = append(results, addr)
		}
	}

	input.KubeletAddresses = results
}
