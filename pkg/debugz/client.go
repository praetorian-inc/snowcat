package debugz

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

type DebugzClient interface {
	Resources(ctx context.Context) ([]runtime.Object, error)
	Version(ctx context.Context) (string, error)
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
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	url := fmt.Sprintf("http://%s/debug/configz", c.debugAddr)
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	log.WithFields(log.Fields{
		"method": req.Method,
		"url":    req.URL.String(),
	}).Debug("validating debug API with HTTP request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code %d from %s", resp.StatusCode, url)
	}
	return nil
}

// Resources queries the Istio debug server for all resources.
func (c *debugzClient) Resources(ctx context.Context) ([]runtime.Object, error) {
	var configs []json.RawMessage

	url := fmt.Sprintf("http://%s/debug/configz", c.debugAddr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	log.WithFields(log.Fields{
		"method": req.Method,
		"url":    req.URL.String(),
	}).Debug("sending HTTP request to debug API")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

func getVersionFromBody(body []byte) (string, error) {
	r := regexp.MustCompile(`istio_version\": \"(.*)\",`)
	matches := r.FindAllSubmatch(body, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("could not find istio_version in syncz debug endpoint")
	}
	return string(matches[0][1]), nil
}

func (c *debugzClient) Version(ctx context.Context) (string, error) {
	url := fmt.Sprintf("http://%s/debug/syncz", c.debugAddr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)

	log.WithFields(log.Fields{
		"method": req.Method,
		"url":    req.URL.String(),
	}).Debug("sending HTTP request to debug API")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return getVersionFromBody(body)
}
