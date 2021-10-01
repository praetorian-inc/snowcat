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
		log.Printf("running %s runner", r.Name)

		err := r.Run(disco)
		if err != nil {
			log.Printf("failed to run %s: %s", r.Name, err)
		}
	}

	if disco.DiscoveryAddress != "" {
		log.Printf("querying xds at %s", disco.DiscoveryAddress)
		cli, err := xds.NewClient(disco.DiscoveryAddress)
		if err != nil {
			log.Printf("failed to initialize xds client: %s", err)
		}
		res, err := cli.Resources(ctx)
		if err != nil {
			log.Printf("failed to query xds resources: %s", err)
		}
		resources.Load(res)
		disco.IstioVersion, err = cli.Version(ctx)
		if err != nil {
			log.Printf("failed to query versions via xds resources: %s", err)
		}
		cli.Close()
	}
	if disco.DebugzAddress != "" {
		log.Printf("querying debug API at %s", disco.DebugzAddress)
		cli, err := debugz.NewClient(disco.DebugzAddress)
		if err != nil {
			log.Printf("failed to initialize debugz client: %s", err)
		}
		res, err := cli.Resources(ctx)
		if err != nil {
			log.Printf("failed to query debugz resources: %s", err)
		}
		disco.IstioVersion, err = cli.Version(ctx)
		if err != nil {
			log.Printf("failed to query versions via debugz resources: %s", err)
		}
		resources.Load(res)
	}
	if len(disco.KubeletAddresses) > 0 {
		for _, addr := range disco.KubeletAddresses {
			cli, err := kubeletclient.NewClient(addr)
			if err != nil {
				log.Printf("failed to create kubelet client: %s", err)
				continue
			}
			pods, err := cli.Pods(ctx)
			if err != nil {
				log.Printf("failed to query pods from kubelet: %s", err)
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
