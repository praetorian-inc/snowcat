package context

import (
	"log"

	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	securityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Resources struct {
	Namespaces            []corev1.Namespace
	PeerAuthentications   []securityv1beta1.PeerAuthentication
	AuthorizationPolicies []securityv1beta1.AuthorizationPolicy
	DestinationRules      []networkingv1alpha3.DestinationRule
	Gateways              []networkingv1alpha3.Gateway
	VirtualServices       []networkingv1alpha3.VirtualService
}

func (r *Resources) Load(resources []runtime.Object) error {
	for _, resource := range resources {
		switch obj := resource.(type) {
		case *securityv1beta1.PeerAuthentication:
			r.PeerAuthentications = append(r.PeerAuthentications, *obj)
		case *securityv1beta1.AuthorizationPolicy:
			r.AuthorizationPolicies = append(r.AuthorizationPolicies, *obj)
		case *networkingv1alpha3.DestinationRule:
			r.DestinationRules = append(r.DestinationRules, *obj)
		case *networkingv1alpha3.Gateway:
			r.Gateways = append(r.Gateways, *obj)
		case *networkingv1alpha3.VirtualService:
			r.VirtualServices = append(r.VirtualServices, *obj)
		case *corev1.Namespace:
			r.Namespaces = append(r.Namespaces, *obj)
		default:
			log.Printf("unknown resource %T", obj)
		}
	}
	return nil
}
