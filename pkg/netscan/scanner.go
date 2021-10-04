package netscan

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type Scanner interface {
	Scan(addr string, timeout time.Duration) bool
}
type Mode string

const (
	ModeHTTP = "http"
	ModeTCP  = "tcp"
)

type ScanJob struct {
	scanner Scanner
	mode    Mode
	hosts   []string
	ports   []string
	lock    *semaphore.Weighted
}

type Option func(*ScanJob)

func WithTLS(config *tls.Config) Option {
	return func(s *ScanJob) {
		if s.mode == "http" {
			hs := s.scanner.(*httpScanner)
			hs.tls = true
		}
	}
}

func New(mode Mode, hosts []string, ports []string, opts ...Option) (*ScanJob, error) {
	var scanner Scanner
	switch mode {
	case ModeHTTP:
		scanner = &httpScanner{}
	case ModeTCP:
		scanner = &tcpScanner{}
	default:
		return nil, fmt.Errorf("unknown scanner mode: %s", mode)
	}
	s := &ScanJob{
		lock:    semaphore.NewWeighted(256),
		scanner: scanner,
		mode:    mode,
		hosts:   hosts,
		ports:   ports,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *ScanJob) Scan(timeout time.Duration) chan string {
	results := make(chan string, 256)

	go func() {
		wg := sync.WaitGroup{}

		defer func() {
			wg.Wait()
			close(results)
		}()

		for _, host := range s.hosts {
			for _, port := range s.ports {
				s.lock.Acquire(context.TODO(), 1) // nolint:errcheck // Acquire will only throw errors from the context
				wg.Add(1)

				addr := net.JoinHostPort(host, port)
				go func(addr string) {
					defer s.lock.Release(1)
					defer wg.Done()

					if s.scanner.Scan(addr, timeout) {
						results <- addr
					}
				}(addr)
			}
		}
	}()

	return results
}
