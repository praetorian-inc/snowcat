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
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	insecureClient = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // nolint:gosec // Scan even if the TLS cert is invalid.
			},
		},
	}
)

type httpScanner struct {
	tls bool
}

func (s *httpScanner) Scan(addr string, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	u := &url.URL{
		Scheme: "http",
		Host:   addr,
	}
	if s.tls {
		u.Scheme = "https"
	}

	req, err := http.NewRequest("HEAD", u.String(), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"url": u.String(),
			"err": err,
		}).Warn("invalid url for port scanning")
	}

	log.WithFields(log.Fields{
		"method": req.Method,
		"url":    req.URL.String(),
	}).Trace("sending HTTP request for port scan")

	resp, err := insecureClient.Do(req.WithContext(ctx))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}
