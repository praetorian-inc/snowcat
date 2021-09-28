package debugz

import (
	"context"
	"fmt"

	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	securityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	operatorv1alpha1 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"

	istiocontext "github.com/praetorian-inc/mithril/pkg/context"
	"github.com/praetorian-inc/mithril/pkg/debugz"
	"github.com/praetorian-inc/mithril/pkg/types"
)

// DebugzContext provides an IstioContext from a directory of static YAML files.
type DebugzContext struct {
	client debugz.DebugzClient

	resources *istiocontext.Resources
}

// New creates a new static IstioContext using yaml files from the provided directory.
func New(addr string) (types.IstioContext, error) {
	err := istioscheme.AddToScheme(clientsetscheme.Scheme)
	if err != nil {
		return nil, err
	}
	client, err := debugz.NewClient(addr)
	if err != nil {
		return nil, err
	}

	c := &DebugzContext{
		client:    client,
		resources: &istiocontext.Resources{},
	}
	return c, c.load()
}

func (c *DebugzContext) load() error {
	ctx := context.Background()

	resources, err := c.client.Resources(ctx)
	if err != nil {
		return err
	}

	if err := c.resources.Load(resources); err != nil {
		return err
	}

	return nil
}

func (c *DebugzContext) IstioNamespace() (string, error) {
	return "", fmt.Errorf("IstioNamespace() unimplemented in static context")
}

func (c *DebugzContext) Namespaces() ([]string, error) {
	var names []string
	for _, ns := range c.resources.Namespaces {
		names = append(names, ns.Name)
	}
	return names, nil
}

func (c *DebugzContext) Version() (string, error) {
	return "", fmt.Errorf("Version() unimplemented in static context")
}

func (c *DebugzContext) IstioOperator() (operatorv1alpha1.IstioOperator, error) {
	return operatorv1alpha1.IstioOperator{}, fmt.Errorf("IstioOperator() unimplemented in static context")
}

func (c *DebugzContext) PeerAuthentications() ([]securityv1beta1.PeerAuthentication, error) {
	return c.resources.PeerAuthentications, nil
}

func (c *DebugzContext) AuthorizationPolicies() ([]securityv1beta1.AuthorizationPolicy, error) {
	return c.resources.AuthorizationPolicies, nil
}

func (c *DebugzContext) DestinationRules() ([]networkingv1alpha3.DestinationRule, error) {
	return c.resources.DestinationRules, nil
}

func (c *DebugzContext) Gateways() ([]networkingv1alpha3.Gateway, error) {
	return c.resources.Gateways, nil
}

func (c *DebugzContext) VirtualServices() ([]networkingv1alpha3.VirtualService, error) {
	return c.resources.VirtualServices, nil
}
