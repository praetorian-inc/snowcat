package types

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	networking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	security "istio.io/client-go/pkg/apis/security/v1beta1"
	securityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	operator "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"

	"github.com/hashicorp/go-multierror"
	"github.com/praetorian-inc/mithril/pkg/util/namer"
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
	Resource    string   `json:"resource"`
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
}

type ObjectGetter interface {
	io.Closer

	Resources(ctx context.Context) []runtime.Object
}

type Resources struct {
	counter int
	decoder runtime.Decoder
	seen    map[string]struct{}

	Namespaces            []corev1.Namespace
	Pods                  []corev1.Pod
	PeerAuthentications   []securityv1beta1.PeerAuthentication
	AuthorizationPolicies []securityv1beta1.AuthorizationPolicy
	DestinationRules      []networkingv1alpha3.DestinationRule
	Gateways              []networkingv1alpha3.Gateway
	VirtualServices       []networkingv1alpha3.VirtualService
	EnvoyFilters          []networkingv1alpha3.EnvoyFilter
	ServiceEntries        []networkingv1alpha3.ServiceEntry
}

func NewResources() Resources {
	istioscheme.AddToScheme(clientsetscheme.Scheme)
	return Resources{
		decoder: clientsetscheme.Codecs.UniversalDeserializer(),
		seen:    make(map[string]struct{}),
	}
}

func (r *Resources) addIfNotExists(obj runtime.Object, meta metav1.ObjectMeta, add func()) {
	gk := obj.GetObjectKind().GroupVersionKind().GroupKind()

	var key string
	if meta.Namespace == "" {
		key = fmt.Sprintf("%s:%s:%s", gk.Group, gk.Kind, meta.Name)
	} else {
		key = fmt.Sprintf("%s:%s:%s:%s", gk.Group, gk.Kind, meta.Namespace, meta.Name)
	}

	if _, ok := r.seen[key]; ok {
		return
	}

	// NOTE: this creates namespaces as we observe them, but without any labels or annotations
	if meta.Namespace != "" {
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: meta.Namespace}}
		r.addIfNotExists(ns, ns.ObjectMeta, func() { r.Namespaces = append(r.Namespaces, *ns) })
	}

	add()
	r.seen[key] = struct{}{}
	r.counter++
}

func (r *Resources) Load(resources []runtime.Object) error {
	for _, resource := range resources {
		switch obj := resource.(type) {
		case *securityv1beta1.PeerAuthentication:
			r.addIfNotExists(resource, obj.ObjectMeta, func() {
				r.PeerAuthentications = append(r.PeerAuthentications, *obj)
			})
		case *securityv1beta1.AuthorizationPolicy:
			r.addIfNotExists(resource, obj.ObjectMeta, func() {
				r.AuthorizationPolicies = append(r.AuthorizationPolicies, *obj)
			})
		case *networkingv1alpha3.DestinationRule:
			r.addIfNotExists(resource, obj.ObjectMeta, func() {
				r.DestinationRules = append(r.DestinationRules, *obj)
			})
		case *networkingv1alpha3.Gateway:
			r.addIfNotExists(resource, obj.ObjectMeta, func() {
				r.Gateways = append(r.Gateways, *obj)
			})
		case *networkingv1alpha3.EnvoyFilter:
			r.addIfNotExists(resource, obj.ObjectMeta, func() {
				r.EnvoyFilters = append(r.EnvoyFilters, *obj)
			})
		case *networkingv1alpha3.VirtualService:
			r.addIfNotExists(resource, obj.ObjectMeta, func() {
				r.VirtualServices = append(r.VirtualServices, *obj)
			})
		case *networkingv1alpha3.ServiceEntry:
			r.addIfNotExists(resource, obj.ObjectMeta, func() {
				r.ServiceEntries = append(r.ServiceEntries, *obj)
			})
		case *corev1.Pod:
			r.addIfNotExists(resource, obj.ObjectMeta, func() {
				r.Pods = append(r.Pods, *obj)
			})
		case *corev1.Namespace:
			r.addIfNotExists(resource, obj.ObjectMeta, func() {
				r.Namespaces = append(r.Namespaces, *obj)
			})
		default:
			log.Printf("unknown resource %T", obj)
		}
	}
	return nil
}

func (r *Resources) LoadFromDirectory(dir string) error {
	root := os.DirFS(dir)

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
		return r.load(data)
	})
}

func (r *Resources) load(data []byte) error {
	var resources []runtime.Object
	for _, yaml := range bytes.Split(data, []byte("---")) {
		if len(yaml) == 0 {
			continue
		}
		obj, _, err := r.decoder.Decode(yaml, nil, nil)
		if err != nil {
			return err
		}

		resources = append(resources, obj)
	}
	return r.Load(resources)
}

func (r *Resources) Len() int {
	return r.counter
}

func (r *Resources) Export(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	var errs error
	lists := []runtime.Object{
		&corev1.NamespaceList{Items: r.Namespaces},
		&corev1.PodList{Items: r.Pods},
		&networkingv1alpha3.DestinationRuleList{Items: r.DestinationRules},
		&networkingv1alpha3.EnvoyFilterList{Items: r.EnvoyFilters},
		&networkingv1alpha3.GatewayList{Items: r.Gateways},
		&networkingv1alpha3.ServiceEntryList{Items: r.ServiceEntries},
		&networkingv1alpha3.VirtualServiceList{Items: r.VirtualServices},
		&securityv1beta1.AuthorizationPolicyList{Items: r.AuthorizationPolicies},
		&securityv1beta1.PeerAuthenticationList{Items: r.PeerAuthentications},
	}
	for _, list := range lists {
		err := exportObjects(dir, list)
		if err != nil {
			errs = multierror.Append(err)
		}
	}
	return errs
}

func exportObjects(dir string, obj runtime.Object) (err error) {
	const mediaType = runtime.ContentTypeYAML
	info, ok := runtime.SerializerInfoForMediaType(clientsetscheme.Codecs.SupportedMediaTypes(), mediaType)
	if !ok {
		return fmt.Errorf("unable to locate encoder -- %q is not a supported media type", mediaType)
	}

	gvks, _, err := clientsetscheme.Scheme.ObjectKinds(obj)
	if err != nil {
		return err
	}
	if len(gvks) == 0 {
		return fmt.Errorf("unknown gvk for %T", obj)
	}

	encoder := clientsetscheme.Codecs.EncoderForVersion(info.Serializer, gvks[0].GroupVersion())

	name := strings.TrimSuffix(gvks[0].Kind, "List")
	filename := filepath.Join(dir, namer.PluralName(name)+".yaml")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		ferr := f.Close()
		if ferr != nil && err == nil {
			err = ferr
		}
	}()

	if err = encoder.Encode(obj, f); err != nil {
		return err
	}

	return nil
}
