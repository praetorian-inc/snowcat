// Copyright 2021 Praetorian Security, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Client wraps methods exposed by the kubelet read-only API.
type Client struct {
	kubeletAddr string

	decoder runtime.Decoder
}

// NewClient creates a kubelet client on the read-only port.
func NewClient(addr string) (*Client, error) {
	cli := &Client{
		kubeletAddr: addr,
		decoder:     clientsetscheme.Codecs.UniversalDeserializer(),
	}
	return cli, cli.verify()
}

func (c *Client) verify() error {
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
func (c *Client) Pods(ctx context.Context) ([]v1.Pod, error) {
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
