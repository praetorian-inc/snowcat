package runner

import "fmt"

const (
	IstioNamespaceKey   = "istioNamespace"
	DiscoveryAddressKey = "discoveryAddress"
)

type Runner struct {
	Strategies []Strategy
}

func (r *Runner) Run(input map[string]string) (map[string]string, error) {
	for _, strategy := range r.Strategies {
		res, err := strategy.Run(input)
		if err != nil {
			return res, nil
		}
	}
	return map[string]string{}, fmt.Errorf("all strategies failed")
}

type Strategy interface {
	Run(map[string]string) (map[string]string, error)
}
