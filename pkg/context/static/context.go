package static

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"

	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	securityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	operatorv1alpha1 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"

	"github.com/praetorian-inc/mithril/pkg/context"
	"github.com/praetorian-inc/mithril/pkg/types"
)

// StaticIstioContext provides an IstioContext from a directory of static YAML files.
type StaticIstioContext struct {
	root    string
	decoder runtime.Decoder

	resources *context.Resources
}

// New creates a new static IstioContext using yaml files from the provided directory.
func New(directory string) (types.IstioContext, error) {
	err := istioscheme.AddToScheme(clientsetscheme.Scheme)
	if err != nil {
		return nil, err
	}

	ctx := &StaticIstioContext{
		root:      directory,
		decoder:   clientsetscheme.Codecs.UniversalDeserializer(),
		resources: &context.Resources{},
	}
	return ctx, ctx.loadAll()
}

func (ctx *StaticIstioContext) loadAll() error {
	root := os.DirFS(ctx.root)

	return fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(root, path)
		if err != nil {
			return err
		}
		return ctx.load(data)
	})
}

func (ctx *StaticIstioContext) load(data []byte) error {
	var resources []runtime.Object
	for _, yaml := range bytes.Split(data, []byte("---")) {
		if len(yaml) == 0 {
			continue
		}
		obj, _, err := ctx.decoder.Decode(yaml, nil, nil)
		if err != nil {
			return err
		}

		resources = append(resources, obj)
	}
	return ctx.resources.Load(resources)
}

func (ctx *StaticIstioContext) IstioNamespace() (string, error) {
	return "", fmt.Errorf("IstioNamespace() unimplemented in static context")
}

func (ctx *StaticIstioContext) Namespaces() ([]string, error) {
	var names []string
	for _, ns := range ctx.resources.Namespaces {
		names = append(names, ns.Name)
	}
	return names, nil
}

func (ctx *StaticIstioContext) Version() (string, error) {
	return "", fmt.Errorf("Version() unimplemented in static context")
}

func (ctx *StaticIstioContext) IstioOperator() (operatorv1alpha1.IstioOperator, error) {
	return operatorv1alpha1.IstioOperator{}, fmt.Errorf("IstioOperator() unimplemented in static context")
}

func (ctx *StaticIstioContext) PeerAuthentications() ([]securityv1beta1.PeerAuthentication, error) {
	return ctx.resources.PeerAuthentications, nil
}

func (ctx *StaticIstioContext) AuthorizationPolicies() ([]securityv1beta1.AuthorizationPolicy, error) {
	return ctx.resources.AuthorizationPolicies, nil
}

func (ctx *StaticIstioContext) DestinationRules() ([]networkingv1alpha3.DestinationRule, error) {
	return ctx.resources.DestinationRules, nil
}

func (ctx *StaticIstioContext) Gateways() ([]networkingv1alpha3.Gateway, error) {
	return ctx.resources.Gateways, nil
}

func (ctx *StaticIstioContext) VirtualServices() ([]networkingv1alpha3.VirtualService, error) {
	return ctx.resources.VirtualServices, nil
}
