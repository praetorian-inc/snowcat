package kubelet

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

type KubeletClient interface {
	Pods(ctx context.Context) ([]v1.Pod, error)
}

type kubeletClient struct {
	kubeletAddr string

	decoder runtime.Decoder
}

// NewClient creates a kubelet client on the read-only port.
func NewClient(addr string) (KubeletClient, error) {
	cli := &kubeletClient{
		kubeletAddr: addr,
		decoder:     clientsetscheme.Codecs.UniversalDeserializer(),
	}
	return cli, cli.verify()
}

func (c *kubeletClient) verify() error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	url := fmt.Sprintf("http://%s/healthz/ping", c.kubeletAddr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

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

// Pods queries the read-only kubelet API for a list of pods running on that node.
func (c *kubeletClient) Pods(ctx context.Context) ([]v1.Pod, error) {
	url := fmt.Sprintf("http://%s/pods", c.kubeletAddr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	obj, _, err := c.decoder.Decode(buf, nil, nil)
	if err != nil {
		return nil, err
	}
	pods, ok := obj.(*v1.PodList)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T from kubelet", obj)
	}
	return pods.Items, nil
}
