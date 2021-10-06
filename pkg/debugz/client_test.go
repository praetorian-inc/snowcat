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

package debugz

import (
	"testing"

	"github.com/bmizerany/assert"
)

const DebugSynczContent = `
[
    {
        "proxy": "istio-ingressgateway-8f568d595-6999c.istio-system",
        "istio_version": "1.10.3",
        "cluster_sent": "lp6Rifqoamc=95893767-fe4c-42a9-96af-d29834fc5e98",
        "cluster_acked": "lp6Rifqoamc=95893767-fe4c-42a9-96af-d29834fc5e98",
        "listener_sent": "lp6Rifqoamc=776d0b13-6b53-4a29-b2a4-f4a8649d03e7",
        "listener_acked": "lp6Rifqoamc=776d0b13-6b53-4a29-b2a4-f4a8649d03e7",
        "route_sent": "lp6Rifqoamc=ef20e494-04e7-46a4-9551-626888b0d09e",
        "route_acked": "lp6Rifqoamc=ef20e494-04e7-46a4-9551-626888b0d09e",
        "endpoint_sent": "lp6Rifqoamc=f387a9ff-8641-448c-b34b-984a816e8a07",
        "endpoint_acked": "lp6Rifqoamc=f387a9ff-8641-448c-b34b-984a816e8a07"
    },
    {
        "proxy": "istio-egressgateway-5547fcc8fc-5bbtd.istio-system",
        "istio_version": "1.10.3",
        "cluster_sent": "lp6Rifqoamc=4d31aadb-332b-466f-82bd-4d7baac559bc",
        "cluster_acked": "lp6Rifqoamc=4d31aadb-332b-466f-82bd-4d7baac559bc",
        "listener_sent": "lp6Rifqoamc=beeed190-43fb-403b-b1d5-05d56d8be77f",
        "listener_acked": "lp6Rifqoamc=beeed190-43fb-403b-b1d5-05d56d8be77f",
        "endpoint_sent": "lp6Rifqoamc=7982c80b-dd7b-4e60-a983-28c973e21900",
        "endpoint_acked": "lp6Rifqoamc=7982c80b-dd7b-4e60-a983-28c973e21900"
    },
    {
        "proxy": "secret-inject-deployment-5754d9b66-cftnb.secret-inject",
        "istio_version": "1.10.3",
        "cluster_sent": "lp6Rifqoamc=b205e809-4256-4810-98ff-4c83fee32f37",
        "cluster_acked": "lp6Rifqoamc=b205e809-4256-4810-98ff-4c83fee32f37",
        "listener_sent": "lp6Rifqoamc=4fdea8fc-856d-4147-b906-f9d03cb9397e",
        "listener_acked": "lp6Rifqoamc=4fdea8fc-856d-4147-b906-f9d03cb9397e",
        "route_sent": "lp6Rifqoamc=7caca542-9a0a-4c70-ae52-9d39f7377b43",
        "route_acked": "lp6Rifqoamc=7caca542-9a0a-4c70-ae52-9d39f7377b43",
        "endpoint_sent": "lp6Rifqoamc=e8ba4216-8341-48a2-b55f-f44f47cf8bea",
        "endpoint_acked": "lp6Rifqoamc=e8ba4216-8341-48a2-b55f-f44f47cf8bea"
    },
    {
        "proxy": "httpbin2-5547847fc4-2dvqx.default",
        "istio_version": "1.10.3",
        "cluster_sent": "lp6Rifqoamc=733118bc-f87c-4826-ba13-8ae6c26b541c",
        "cluster_acked": "lp6Rifqoamc=733118bc-f87c-4826-ba13-8ae6c26b541c",
        "listener_sent": "lp6Rifqoamc=ef4b1961-24e8-4e8c-9cb1-f5f2c53b238f",
        "listener_acked": "lp6Rifqoamc=ef4b1961-24e8-4e8c-9cb1-f5f2c53b238f",
        "route_sent": "lp6Rifqoamc=6f336972-389c-4c29-a525-eccbf6ce36d9",
        "route_acked": "lp6Rifqoamc=6f336972-389c-4c29-a525-eccbf6ce36d9",
        "endpoint_sent": "lp6Rifqoamc=6c97ba82-9efd-404f-bacd-6b7c9eff12ac",
        "endpoint_acked": "lp6Rifqoamc=6c97ba82-9efd-404f-bacd-6b7c9eff12ac"
    },
    {
        "proxy": "httpbin1-5848b579fb-fhd4j.default",
        "istio_version": "1.10.3",
        "cluster_sent": "lp6Rifqoamc=66dba551-4357-426b-a024-d90c5c3c7214",
        "cluster_acked": "lp6Rifqoamc=66dba551-4357-426b-a024-d90c5c3c7214",
        "listener_sent": "lp6Rifqoamc=deae7179-0cc9-47a8-93aa-1e28619a1e7d",
        "listener_acked": "lp6Rifqoamc=deae7179-0cc9-47a8-93aa-1e28619a1e7d",
        "route_sent": "lp6Rifqoamc=ef135cbb-2f9a-48a0-bd1e-f0a085a6437f",
        "route_acked": "lp6Rifqoamc=ef135cbb-2f9a-48a0-bd1e-f0a085a6437f",
        "endpoint_sent": "lp6Rifqoamc=51aac419-962c-478b-a4f4-69cbdb2e7d4e",
        "endpoint_acked": "lp6Rifqoamc=51aac419-962c-478b-a4f4-69cbdb2e7d4e"
    }
]`

func TestGetMatchingVulns(t *testing.T) {
	version, _ := getVersionFromBody(([]byte)(DebugSynczContent))
	assert.Equal(t, "1.10.3", version)
}

// only usable if localhost 8080 is hosting this endpoint
// func TestGetMatchingVulnsFromEndpoint(t *testing.T) {
// 	url := "http://127.0.0.1:8080/debug/syncz"
// 	req, _ := http.NewRequest("GET", url, nil)
// 	resp, _ := http.DefaultClient.Do(req)
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	version, _ := getVersionFromBody(body)
// 	assert.Equal(t, "1.10.3", version)
// }
