package types

import (
	networking "istio.io/client-go/pkg/apis/networking/v1beta1"
	security "istio.io/client-go/pkg/apis/security/v1beta1"
)

type Config struct {
}

type IstioContext interface {
	IstioNamespace() (string, error)
	Namespaces() ([]string, error)
	Version() (string, error)
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
	Audit(c IstioContext) ([]AuditResult, error)
}
