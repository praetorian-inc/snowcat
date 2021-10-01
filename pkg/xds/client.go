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

package xds

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/gogo/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	mcp "istio.io/api/mcp/v1alpha1"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	"istio.io/pkg/version"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"

	blockinggrpc "github.com/praetorian-inc/mithril/pkg/grpc"
)

type DiscoveryClient interface {
	Version(ctx context.Context) (string, error)
	List(ctx context.Context, gvk schema.GroupVersionKind) ([]runtime.Object, error)
	Resources(ctx context.Context) ([]runtime.Object, error)
	Close() error
}

type xdsClient struct {
	discoveryAddr string
	opts          []grpc.DialOption

	conn   *grpc.ClientConn
	connMu sync.Mutex

	stream discovery.AggregatedDiscoveryService_StreamAggregatedResourcesClient

	decoder runtime.Decoder
}

// NewClient creates an XDS client given a GRPC address.
func NewClient(addr string) (DiscoveryClient, error) {
	err := istioscheme.AddToScheme(clientsetscheme.Scheme)
	if err != nil {
		return nil, err
	}

	cli := &xdsClient{
		discoveryAddr: addr,
		opts: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		decoder: clientsetscheme.Codecs.UniversalDeserializer(),
	}
	_, err = cli.Version(context.Background())
	return cli, err
}

func (xds *xdsClient) makeNodeID() string {
	// TODO: should we attempt to populate this?
	return "sidecar~0.0.0.0~mithril~mithril"
}

func (xds *xdsClient) makeRequest(typeURL string) *discovery.DiscoveryRequest {
	return &discovery.DiscoveryRequest{
		Node: &core.Node{
			Id: xds.makeNodeID(),
		},
		TypeUrl: typeURL,
	}
}

func (xds *xdsClient) connect(ctx context.Context) error {
	xds.connMu.Lock()
	defer xds.connMu.Unlock()

	var err error

	connctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	log.WithFields(log.Fields{
		"addr": xds.discoveryAddr,
	}).Debug("connecting to xds")

	xds.conn, err = blockinggrpc.BlockingDial(connctx, "tcp", xds.discoveryAddr, nil, xds.opts...)
	if err != nil {
		return err
	}

	xds.stream, err = discovery.NewAggregatedDiscoveryServiceClient(xds.conn).
		StreamAggregatedResources(ctx)
	if err != nil {
		xds.Close()
		return err
	}

	return nil
}

func (xds *xdsClient) Close() error {
	xds.connMu.Lock()
	defer xds.connMu.Unlock()

	if xds.conn != nil {
		err := xds.conn.Close()
		xds.conn = nil
		xds.stream = nil
		return err
	}
	return nil
}

func (xds *xdsClient) send(ctx context.Context, req *discovery.DiscoveryRequest) (*discovery.DiscoveryResponse, error) {
	if xds.conn == nil {
		err := xds.connect(ctx)
		if err != nil {
			return nil, err
		}
	}

	log.WithFields(log.Fields{
		"typeURL": req.TypeUrl,
	}).Debug("sending xds request")

	err := xds.stream.Send(req)
	if err != nil {
		return nil, err
	}
	return xds.stream.Recv()
}

// Version queries the XDS server and retrieves its version.
// {
//   "Component": "istiod",
//   "ID": "istiod-568d797f55-vlxmt",
//   "Info": {
//     "version": "1.10.3",
//     "revision": "61313778e0b785e401c696f5e92f47af069f96d0",
//     "golang_version": "go1.16.6",
//     "status": "Clean",
//     "tag": "1.10.3"
//   }
// }
func (xds *xdsClient) Version(ctx context.Context) (string, error) {
	req := xds.makeRequest("")
	res, err := xds.send(ctx, req)
	if err != nil {
		return "", err
	}

	ident := res.ControlPlane.Identifier

	var info version.ServerInfo
	err = json.Unmarshal([]byte(ident), &info)
	if err != nil {
		return "", err
	}
	return info.Info.Version, nil
}

func getObjectMetadata(m *mcp.Resource) (metav1.ObjectMeta, error) {
	if m == nil || m.Metadata == nil {
		return metav1.ObjectMeta{}, nil
	}
	meta := metav1.ObjectMeta{
		ResourceVersion: m.Metadata.Version,
		Labels:          m.Metadata.Labels,
		Annotations:     m.Metadata.Annotations,
	}
	nsn := strings.Split(m.Metadata.Name, "/")
	if len(nsn) != 2 {
		return metav1.ObjectMeta{}, fmt.Errorf("invalid name %s", m.Metadata.Name)
	}
	meta.Namespace = nsn[0]
	meta.Name = nsn[1]

	return meta, nil
}

func decodeMCPResource(data []byte, gvk schema.GroupVersionKind) (runtime.Object, error) {
	r := &mcp.Resource{}
	if err := proto.Unmarshal(data, r); err != nil {
		return nil, err
	}

	obj, err := istioscheme.Scheme.New(gvk)
	if err != nil {
		return nil, err
	}

	objVal := reflect.ValueOf(obj)
	if objVal.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("object %T is not a Ptr", obj)
	}
	metaVal := objVal.Elem().FieldByName("ObjectMeta")
	if !metaVal.CanSet() {
		return nil, fmt.Errorf("object %T cannot set ObjectMeta", obj)
	}
	meta, ok := metaVal.Addr().Interface().(*metav1.ObjectMeta)
	if !ok {
		return nil, fmt.Errorf("object %s not of type metav1.ObjectMeta", metaVal.Kind())
	}
	*meta, err = getObjectMetadata(r)
	if err != nil {
		return nil, err
	}

	if objVal.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("object %T is not a Ptr", obj)
	}
	specVal := objVal.Elem().FieldByName("Spec")
	if !specVal.CanSet() {
		return nil, fmt.Errorf("object %T cannot set Spec", obj)
	}
	spec := specVal.Addr().Interface()

	pb, ok := spec.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("object %T does not implement the protobuf marshalling interface", spec)
	}
	if err := proto.Unmarshal(r.Body.Value, pb); err != nil {
		return nil, err
	}

	obj.GetObjectKind().SetGroupVersionKind(gvk)
	return obj, nil
}

// Resources queries the XDS server for a given GroupVersionKind
// (e.g. security.istio.io/v1beta1/AuthorizationPolicy) and
// returns these resources as Kubernetes runtime.Objects.
func (xds *xdsClient) List(ctx context.Context, gvk schema.GroupVersionKind) ([]runtime.Object, error) {
	typeURL := fmt.Sprintf("%s/%s/%s", gvk.Group, gvk.Version, gvk.Kind)
	req := xds.makeRequest(typeURL)
	resp, err := xds.send(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.TypeUrl != typeURL {
		return nil, fmt.Errorf("unexpected typeUrl: %s", resp.TypeUrl)
	}

	var resources []runtime.Object
	for _, res := range resp.Resources {
		obj, err := decodeMCPResource(res.Value, gvk)
		if err != nil {
			return nil, err
		}
		resources = append(resources, obj)
	}

	return resources, nil
}

func (xds *xdsClient) Resources(ctx context.Context) ([]runtime.Object, error) {
	var resources []runtime.Object
	for gvk := range istioscheme.Scheme.AllKnownTypes() {
		res, err := xds.List(ctx, gvk)
		if err != nil {
			return nil, err
		}
		resources = append(resources, res...)
	}
	return resources, nil
}
