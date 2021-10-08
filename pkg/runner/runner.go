// Copyright 2021 Praetorian Security, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package runner provides an abstraction layer: Strategy, for general
// collection of information about a cluster and its components. Additionally,
// it defines a struct: Runner that is comprised of a list of strategies to
// carry out. this allows for a data collection goal to have multiple methods of
// obtaining the desired data.
//
// strategies take as inputs a types.Discovery, which contains various fields
// useful in the introspection of an istio system/cluster. the interface is
// constructed to allow each strategy the opportunity to modify a single
// Discovery instance (by passing a reference to one).
//
// to perform a collection, a consuming package will need to construct a Runners
// struct, containing the list of individual runners desired in the collection,
// and call the Run() method on it.
package runner

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/praetorian-inc/mesh-hunter/pkg/debugz"
	kubeletclient "github.com/praetorian-inc/mesh-hunter/pkg/kubelet"
	"github.com/praetorian-inc/mesh-hunter/pkg/types"
	"github.com/praetorian-inc/mesh-hunter/pkg/xds"
)

// Runner is a struct containing several strategies to try in cluster enumeration.
type Runner struct {
	Name       string
	Strategies []Strategy
}

// Strategy is an interface that abstractly describes how information is to be
// collected in an istio system.
type Strategy interface {
	// Name merely returns the strategy's name for reporting purposes.
	Name() string
	// Run is called on every strategy to execute its particular method
	// of collection. it is passed a reference to a types.Discovery struct to record
	// gathered and verified data. it is assumed that the implementation of this
	// function will perform validation of discovered data before recording it to
	// the Discovery struct.
	Run(input *types.Discovery) error
}

// Run as defined for a runner loops over all strategies and passes a
// *types.Discovery. it then surfaces any errors it receives. if all strategies
// fail, an error is produced
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

// Runners defines a type alias for a list of Runner structs
type Runners []Runner

// Run as defined for a Runners type iterates over all runners in the list, and
// calls each runner's Run() method. it then looks at the resulting
// types.Discovery and attempts to use them to confirm if they are correct.
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
