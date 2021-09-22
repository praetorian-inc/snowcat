package static

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"

	networkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	securityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	operatorv1alpha1 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"

	"github.com/praetorian-inc/mithril/pkg/types"
)

// StaticIstioContext provides an IstioContext from a directory of static YAML files.
type StaticIstioContext struct {
	root    string
	decoder runtime.Decoder

	namespaces            []corev1.Namespace
	peerAuthentications   []securityv1beta1.PeerAuthentication
	authorizationPolicies []securityv1beta1.AuthorizationPolicy
	destinationRules      []networkingv1beta1.DestinationRule
	gateways              []networkingv1beta1.Gateway
	virtualServices       []networkingv1beta1.VirtualService
}

// New creates a new static IstioContext using yaml files from the provided directory.
func New(directory string) (types.IstioContext, error) {
	err := istioscheme.AddToScheme(clientsetscheme.Scheme)
	if err != nil {
		return nil, err
	}

	ctx := &StaticIstioContext{
		root:    directory,
		decoder: clientsetscheme.Codecs.UniversalDeserializer(),
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
	for _, yaml := range bytes.Split(data, []byte("---")) {
		if len(yaml) == 0 {
			continue
		}
		obj, gvk, err := ctx.decoder.Decode(yaml, nil, nil)
		if err != nil {
			return err
		}
		log.Printf("resource %s", gvk.String())

		switch resource := obj.(type) {
		case *securityv1beta1.PeerAuthentication:
			ctx.peerAuthentications = append(ctx.peerAuthentications, *resource)
		case *securityv1beta1.AuthorizationPolicy:
			ctx.authorizationPolicies = append(ctx.authorizationPolicies, *resource)
		case *networkingv1beta1.DestinationRule:
			ctx.destinationRules = append(ctx.destinationRules, *resource)
		case *networkingv1beta1.Gateway:
			ctx.gateways = append(ctx.gateways, *resource)
		case *networkingv1beta1.VirtualService:
			ctx.virtualServices = append(ctx.virtualServices, *resource)
		case *corev1.Namespace:
			ctx.namespaces = append(ctx.namespaces, *resource)
		default:
			log.Printf("skipping unknown resource %s", gvk.String())
		}
	}
	return nil
}

func (ctx *StaticIstioContext) IstioNamespace() (string, error) {
	return "", fmt.Errorf("IstioNamespace() unimplemented in static context")
}

func (ctx *StaticIstioContext) Namespaces() ([]string, error) {
	var names []string
	for _, ns := range ctx.namespaces {
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
	return ctx.peerAuthentications, nil
}

func (ctx *StaticIstioContext) AuthorizationPolicies() ([]securityv1beta1.AuthorizationPolicy, error) {
	return ctx.authorizationPolicies, nil
}

func (ctx *StaticIstioContext) DestinationRules() ([]networkingv1beta1.DestinationRule, error) {
	return ctx.destinationRules, nil
}

func (ctx *StaticIstioContext) Gateways() ([]networkingv1beta1.Gateway, error) {
	return ctx.gateways, nil
}

func (ctx *StaticIstioContext) VirtualServices() ([]networkingv1beta1.VirtualService, error) {
	return ctx.virtualServices, nil
}
