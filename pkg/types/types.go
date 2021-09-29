package types

import (
	"context"
	"io"
	"log"

	networking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	security "istio.io/client-go/pkg/apis/security/v1beta1"
	securityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	operator "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Config struct {
}

type IstioContext interface {
	IstioNamespace() (string, error)
	Namespaces() ([]string, error)
	Version() (string, error)
	IstioOperator() (operator.IstioOperator, error)
	PeerAuthentications() ([]security.PeerAuthentication, error)
	AuthorizationPolicies() ([]security.AuthorizationPolicy, error)
	DestinationRules() ([]networking.DestinationRule, error)
	Gateways() ([]networking.Gateway, error)
	VirtualServices() ([]networking.VirtualService, error)
}

type AuditResult struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Severity    Severity `json:"severity"`
}

type Severity int

const (
	Unknown Severity = iota
)

type Auditor interface {
	Name() string
	Audit(Discovery, Resources) ([]AuditResult, error)
}

type Discovery struct {
	// IstioVersion is the version of the istio control plane.
	IstioVersion string
	// IstioNamespace is the Kubernetes namespace of the istio control plane.
	IstioNamespace string
	// DiscoveryAddress is the IP:port of istiod's unauthenticated xds.
	DiscoveryAddress string
	// DebugzAddress is the IP:port of istiod's debug API.
	DebugzAddress string
	// KubeletAddresses is a list of addresses of each node's kubelet read-only API.
	// These addresses have the form "host:port".
	KubeletAddresses []string
	// IstiodIPs is a list of IP addresses of that appear to be the istio control plane.
	IstiodIPs []string
}

type ObjectGetter interface {
	io.Closer

	Resources(ctx context.Context) []runtime.Object
}

type Resources struct {
	counter int

	Namespaces            []corev1.Namespace
	PeerAuthentications   []securityv1beta1.PeerAuthentication
	AuthorizationPolicies []securityv1beta1.AuthorizationPolicy
	DestinationRules      []networkingv1alpha3.DestinationRule
	Gateways              []networkingv1alpha3.Gateway
	VirtualServices       []networkingv1alpha3.VirtualService
	Filters               []networkingv1alpha3.EnvoyFilter
}

func (r *Resources) Load(resources []runtime.Object) error {
	for _, resource := range resources {
		switch obj := resource.(type) {
		case *securityv1beta1.PeerAuthentication:
			r.PeerAuthentications = append(r.PeerAuthentications, *obj)
			r.counter++
		case *securityv1beta1.AuthorizationPolicy:
			r.AuthorizationPolicies = append(r.AuthorizationPolicies, *obj)
			r.counter++
		case *networkingv1alpha3.DestinationRule:
			r.DestinationRules = append(r.DestinationRules, *obj)
			r.counter++
		case *networkingv1alpha3.Gateway:
			r.Gateways = append(r.Gateways, *obj)
			r.counter++
		case *networkingv1alpha3.EnvoyFilter:
			r.Filters = append(r.Filters, *obj)
			r.counter++
		case *networkingv1alpha3.VirtualService:
			r.VirtualServices = append(r.VirtualServices, *obj)
			r.counter++
		case *corev1.Namespace:
			r.Namespaces = append(r.Namespaces, *obj)
			r.counter++
		default:
			log.Printf("unknown resource %T", obj)
		}
	}
	return nil
}

func (r *Resources) Len() int {
	return r.counter
}
