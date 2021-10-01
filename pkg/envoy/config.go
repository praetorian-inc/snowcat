package envoy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
)

type EnvoyConfig struct {
	jpathNode *ajson.Node
}

func (ec *EnvoyConfig) DiscoveryAddress() (string, error) {
	// For when the discoveryAddress property exists
	nodes, err := ec.jpathNode.JSONPath("$..discoveryAddress")
	if err == nil && nodes != nil && len(nodes) > 0 {
		addressStringVal := nodes[0].String()
		addressStringVal = strings.Replace(addressStringVal, "\"", "", -1)
		return addressStringVal, nil
	}

	// For when there is a dynamic route config to port 15010 or 15012
	nodes, err = ec.jpathNode.JSONPath("$..name")
	if err == nil && nodes != nil && len(nodes) > 0 {
		for _, node := range nodes {
			addressStringVal := node.String()
			addressStringVal = strings.Replace(addressStringVal, "\"", "", -1)
			if strings.HasSuffix(addressStringVal, ":15010") || strings.HasSuffix(addressStringVal, ":15012") {
				return addressStringVal, nil
			}
		}
	}

	if err != nil {
		return "", err
	}

	return "", fmt.Errorf("Could not find discovery address in Envoy config")
}

func LoadConfig(configBytes []byte) (*EnvoyConfig, error) {
	root, err := ajson.Unmarshal(configBytes)
	if err != nil {
		return nil, err
	}
	return &EnvoyConfig{jpathNode: root}, nil
}

func RetrieveConfig(envoyAdminUrl string) (*EnvoyConfig, error) {
	log.WithFields(log.Fields{
		"method": "GET",
		"url":    envoyAdminUrl,
	}).Debug("sending HTTP request to envoy")

	resp, err := http.Get(envoyAdminUrl)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return LoadConfig(body)
}
