package runner

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/praetorian-inc/mithril/pkg/debugz"
	kubeletclient "github.com/praetorian-inc/mithril/pkg/kubelet"
	"github.com/praetorian-inc/mithril/pkg/types"
	"github.com/praetorian-inc/mithril/pkg/xds"
)

type Runner struct {
	Name       string
	Strategies []Strategy
}

type Strategy interface {
	Name() string
	Run(input *types.Discovery) error
}

func (r *Runner) Run(input *types.Discovery) error {
	var errs error
	for _, strategy := range r.Strategies {
		log.WithFields(log.Fields{
			"runner":   r.Name,
			"strategy": strategy.Name(),
		}).Info("running discovery strategy")
		err := strategy.Run(input)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("%s: %w", strategy.Name(), err))
			continue
		}
		return nil
	}
	return fmt.Errorf("all strategies failed: %s", errs)
}

type Runners []Runner

func (runners Runners) Run(disco *types.Discovery, resources *types.Resources) {
	ctx := context.Background()

	for _, r := range runners {
		err := r.Run(disco)
		if err != nil {
			log.WithFields(log.Fields{
				"runner": r.Name,
				"err":    err,
			}).Warn("failed to run")
		}
	}

	if disco.DiscoveryAddress != "" {
		cli, err := xds.NewClient(disco.DiscoveryAddress)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": disco.DiscoveryAddress,
				"err":  err,
			}).Warn("failed initialize xds client")
		}
		res, err := cli.Resources(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": disco.DiscoveryAddress,
				"err":  err,
			}).Warn("failed query xds resources")
		}
		resources.Load(res)
		disco.IstioVersion, err = cli.Version(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": disco.DiscoveryAddress,
				"err":  err,
			}).Warn("failed query xds version")
		}
		cli.Close()
	}
	if disco.DebugzAddress != "" {
		cli, err := debugz.NewClient(disco.DebugzAddress)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": disco.DebugzAddress,
				"err":  err,
			}).Warn("failed initialize debugz client")
		}
		res, err := cli.Resources(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": disco.DebugzAddress,
				"err":  err,
			}).Warn("failed query debugz resources")
		}
		disco.IstioVersion, err = cli.Version(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": disco.DebugzAddress,
				"err":  err,
			}).Warn("failed query debugz version")
		}
		resources.Load(res)
	}
	if len(disco.KubeletAddresses) > 0 {
		for _, addr := range disco.KubeletAddresses {
			cli, err := kubeletclient.NewClient(addr)
			if err != nil {
				log.WithFields(log.Fields{
					"addr": addr,
					"err":  err,
				}).Warn("failed initialize kubelet client")
				continue
			}
			pods, err := cli.Pods(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"addr": addr,
					"err":  err,
				}).Warn("failed query kubelet pods")
				continue
			}
			var res []runtime.Object
			for i := range pods {
				res = append(res, &pods[i])
			}
			resources.Load(res)
		}
	}
}
