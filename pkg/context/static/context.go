package static

import (
	"fmt"

	networking "istio.io/client-go/pkg/apis/networking/v1beta1"
	security "istio.io/client-go/pkg/apis/security/v1beta1"

	"github.com/praetorian-inc/mithril/pkg/types"
)

// StaticIstioContext provides an IstioContext from a directory of static YAML files.
type StaticIstioContext struct {
	root string
}

// New creates a new static IstioContext using yaml files from the provided directory.
func New(directory string) (types.IstioContext, error) {
	ctx := &StaticIstioContext{
		root: directory,
	}
	return ctx, ctx.load()
}

func (ctx *StaticIstioContext) load() error {
	/*
		decoder := scheme.Codecs.UniversalDeserializer()
		_ = decoder
		for _, yaml := range strings.Split(string(data), "---") {
			if len(yaml) == 0 {
				continue
			}
		}
	*/
	return nil
}

func (ctx *StaticIstioContext) IstioNamespace() (string, error) {
	return "", fmt.Errorf("IstioNamespace() unimplemented in static context")
}

func (ctx *StaticIstioContext) Namespaces() ([]string, error) {
	return nil, fmt.Errorf("Namespaces() unimplemented in static context")
}

func (ctx *StaticIstioContext) Version() (string, error) {
	return "", fmt.Errorf("Version() unimplemented in static context")
}

func (ctx *StaticIstioContext) PeerAuthentications() ([]security.PeerAuthentication, error) {
	return nil, fmt.Errorf("PeerAuthentications() unimplemented in static context")
}

func (ctx *StaticIstioContext) AuthorizationPolicies() ([]security.AuthorizationPolicy, error) {
	return nil, fmt.Errorf("AuthorizationPolicies() unimplemented in static context")
}

func (ctx *StaticIstioContext) DestinationRules() ([]networking.DestinationRule, error) {
	return nil, fmt.Errorf("DestinationRules() unimplemented in static context")
}

func (ctx *StaticIstioContext) Gateways() ([]networking.Gateway, error) {
	return nil, fmt.Errorf("Gateways() unimplemented in static context")
}

func (ctx *StaticIstioContext) VirtualServices() ([]networking.VirtualService, error) {
	return nil, fmt.Errorf("VirtualServices() unimplemented in static context")
}
