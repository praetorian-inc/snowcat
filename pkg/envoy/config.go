package envoy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
)

// Config wraps the Envoy config_dump and exposes methods to extract data from it.
type Config struct {
	jpathNode *ajson.Node
}

// DiscoveryAddress extracts the discoveryAddress property from the config_dump.
func (ec *Config) DiscoveryAddress() (string, error) {
	// For when the discoveryAddress property exists
	nodes, err := ec.jpathNode.JSONPath("$..discoveryAddress")
	if err == nil && nodes != nil && len(nodes) > 0 {
		return strings.ReplaceAll(nodes[0].String(), "\"", ""), nil
	}

	// For when there is a dynamic route config to port 15010 or 15012
	nodes, err = ec.jpathNode.JSONPath("$..name")
	if err != nil {
		return "", err
	}

	for _, node := range nodes {
		addr := strings.ReplaceAll(node.String(), "\"", "")
		if strings.HasSuffix(addr, ":15010") || strings.HasSuffix(addr, ":15012") {
			return addr, nil
		}
	}

	return "", fmt.Errorf("Could not find discovery address in Envoy config")
}

// LoadConfig unmarshals bytes into a Config
func LoadConfig(configBytes []byte) (*Config, error) {
	root, err := ajson.Unmarshal(configBytes)
	if err != nil {
		return nil, err
	}
	return &Config{jpathNode: root}, nil
}

// RetrieveConfig fetches a Config from a local envoy service.
func RetrieveConfig(envoyAdminURL string) (*Config, error) {
	log.WithFields(log.Fields{
		"method": "GET",
		"url":    envoyAdminURL,
	}).Debug("sending HTTP request to envoy")

	resp, err := http.Get(envoyAdminURL) // nolint:gosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return LoadConfig(body)
}
