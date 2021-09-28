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

// NewClient creates a client for the istiod debug API.
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
