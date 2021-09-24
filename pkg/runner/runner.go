package runner

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

const (
	IstioNamespaceKey   = "istioNamespace"
	DiscoveryAddressKey = "discoveryAddress"
)

type Runner struct {
	Name       string
	Strategies []Strategy
}

func (r *Runner) Run(input map[string]string) error {
	var errs error
	for _, strategy := range r.Strategies {
		res, err := strategy.Run(input)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("%s: %w", strategy.Name(), err))
			continue
		}
		err = strategy.Verify(res)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("%s validator: %w", strategy.Name(), err))
			continue
		}

		for k, v := range res {
			input[k] = v
		}
		return nil
	}
	return fmt.Errorf("all strategies failed: %s", errs)
}

type Strategy interface {
	Name() string
	Run(map[string]string) (map[string]string, error)
	Verify(map[string]string) error
}
