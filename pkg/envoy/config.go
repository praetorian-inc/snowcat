package envoy

import (
	"fmt"
	"github.com/spyzhov/ajson"
	"io/ioutil"
	"net/http"
	"strings"
)

type EnvoyConfig struct {
	jpathNode *ajson.Node
}

func (ec *EnvoyConfig) DiscoveryAddress() (string, error) {
	// configs[(first config where @type is type.googleapis.com/envoy.admin.v3.BootstrapConfigDump)]->bootstrap->node->metadata->PROXY_CONFIG->discoveryAddress
	nodes, err := ec.jpathNode.JSONPath("$..discoveryAddress")
	if err != nil {
		return "", err
	}
	if nodes == nil || len(nodes) < 1 {
		return "", fmt.Errorf("no discoveryAddress in Envoy config")
	}
	addressStringVal := nodes[0].String()
	addressStringVal = strings.Replace(addressStringVal, "\"", "", -1)
	return addressStringVal, nil
}

func LoadConfig(configBytes []byte) (*EnvoyConfig, error) {
	root, err := ajson.Unmarshal(configBytes)
	if err != nil {
		return nil, err
	}
	return &EnvoyConfig{jpathNode: root}, nil
}

func RetrieveConfig(envoyAdminUrl string) (*EnvoyConfig, error) {
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
