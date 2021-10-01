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
	InsecureClient = &http.Client{
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
				InsecureSkipVerify: true,
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
	}).Debug("sending HTTP request for port scan")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}
