package types

import (
	networking "istio.io/client-go/pkg/apis/networking/v1beta1"
	security "istio.io/client-go/pkg/apis/security/v1beta1"
)

type Config struct {
}

type IstioContext interface {
	IstioNamespace() string
	Namespaces() []string
	Version() string
	PeerAuthentications() []security.PeerAuthentication
	AuthorizationPolicies() []security.AuthorizationPolicy
	DestinationRules() []networking.DestinationRule
	Gateways() []networking.Gateway
	VirtualServices() []networking.VirtualService
	// IstioOperator() operator.IstioOperatorSpec
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
