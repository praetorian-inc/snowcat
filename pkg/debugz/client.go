// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package debugz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type DebugzClient interface {
	Resources(ctx context.Context) ([]runtime.Object, error)
}

type debugzClient struct {
	debugAddr string

	decoder runtime.Decoder
}

// NewClient creates an XDS client given a GRPC address.
func NewClient(addr string) (DebugzClient, error) {
	err := istioscheme.AddToScheme(clientsetscheme.Scheme)
	if err != nil {
		return nil, err
	}

	cli := &debugzClient{
		debugAddr: addr,
		decoder:   clientsetscheme.Codecs.UniversalDeserializer(),
	}
	return cli, cli.verify()
}

func (c *debugzClient) verify() error {
	url := fmt.Sprintf("http://%s/debug/configz", c.debugAddr)
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code %d from %s", resp.StatusCode, url)
	}
	return nil
}

// Resources queries the Istio debug server for all resources.
func (c *debugzClient) Resources(ctx context.Context) ([]runtime.Object, error) {
	var configs []json.RawMessage

	resp, err := http.Get(fmt.Sprintf("http://%s/debug/configz", c.debugAddr))
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&configs)
	if err != nil {
		return nil, err
	}

	var resources []runtime.Object
	for _, config := range configs {
		obj, _, err := c.decoder.Decode(config, nil, nil)
		if err != nil {
			return nil, err
		}
		resources = append(resources, obj)
	}

	return resources, nil
}
