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

package netscan

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var (
	BadGatewayHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, "Bad Gateway", 502)
		fmt.Fprint(rw, "upstream connect error or disconnect/reset before headers. reset reason: protocol error")
	})
	OkHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, "ok", 200)
	})
)

func TestHTTP(t *testing.T) {
	type testcase struct {
		server *httptest.Server
		open   bool
	}

	testcases := []testcase{
		{
			server: httptest.NewServer(BadGatewayHandler),
			open:   true,
		},
		{
			server: httptest.NewServer(OkHandler),
			open:   true,
		},
		{
			server: &httptest.Server{
				// Assigned as TEST-NET-3, should fail.
				URL: "http://203.0.113.1:1337",
			},
			open: false,
		},
	}

	for i, tc := range testcases {
		u, err := url.Parse(tc.server.URL)
		if err != nil {
			t.Fatalf("[%d] failed to parse httptest server URL: %s", i, err)
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			t.Fatalf("[%d] failed to parse split host %s: %s", i, tc.server.URL, err)
		}

		scanner, err := New(ModeHTTP, []string{host}, []string{port})
		if err != nil {
			t.Fatalf("[%d] failed to init scanner: %s", i, err)
		}
		var results []string
		for addr := range scanner.Scan(500 * time.Millisecond) {
			results = append(results, addr)
		}
		if len(results) > 1 {
			t.Errorf("[%d] unexpected number of results: %d", i, len(results))
		} else if len(results) == 0 && tc.open {
			t.Errorf("[%d] scanner failed to detect open port", i)
		} else if len(results) == 1 && !tc.open {
			t.Errorf("[%d] scanner failed to detect closed port", i)
		}
	}
}

func TestTCP(t *testing.T) {
	type testcase struct {
		addr string
		open bool
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on local port: %s", err)
	}
	defer listener.Close()

	testcases := []testcase{
		{
			addr: listener.Addr().String(),
			open: true,
		},
		{
			addr: "203.0.113.1:1337",
			open: false,
		},
	}

	for i, tc := range testcases {
		host, port, err := net.SplitHostPort(tc.addr)
		if err != nil {
			t.Fatalf("[%d] failed to parse split host %s: %s", i, tc.addr, err)
		}

		scanner, err := New(ModeTCP, []string{host}, []string{port})
		if err != nil {
			t.Fatalf("[%d] failed to init scanner: %s", i, err)
		}
		var results []string
		for addr := range scanner.Scan(500 * time.Millisecond) {
			results = append(results, addr)
		}
		if len(results) > 1 {
			t.Errorf("[%d] unexpected number of results: %d", i, len(results))
		} else if len(results) == 0 && tc.open {
			t.Errorf("[%d] scanner failed to detect open port", i)
		} else if len(results) == 1 && !tc.open {
			t.Errorf("[%d] scanner failed to detect closed port", i)
		}
	}
}
