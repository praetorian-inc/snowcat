package main

import (
	networking "istio.io/api/networking/v1beta1"
	operator "istio.io/api/operator/v1alpha1"
	security "istio.io/api/security/v1beta1"
)

type IstioContext interface {
	IstioNamespace() string
	Namespaces() []string
	Version() string
	PeerAuthentications() []security.PeerAuthentication
	AuthorizationPolicies() []security.AuthorizationPolicy
	DestinationRules() []networking.DestinationRule
	Gateways() []networking.Gateway
	VirtualServices() []networking.VirtualService
	IstioOperator() operator.IstioOperatorSpec
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
	Audit(c IstioContext) ([]AuditResult, error)
}
