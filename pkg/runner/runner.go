package runner

import (
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/praetorian-inc/mithril/pkg/types"
)

type Runner struct {
	Name       string
	Strategies []Strategy
}

func (r *Runner) Run(input *types.Discovery) error {
	var errs error
	for _, strategy := range r.Strategies {
		err := strategy.Run(input)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("%s: %w", strategy.Name(), err))
			continue
		}
		return nil
	}
	return fmt.Errorf("all strategies failed: %s", errs)
}

type Strategy interface {
	Name() string
	Run(input *types.Discovery) error
}
