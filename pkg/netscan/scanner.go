// Package netscan provides functionality for running network enumeration scans.
// currently the scanner provides two modes:
//
// ModeHTTP:
// useful when the services needing enumeration are HTTP servers, and a
// traditional portscan would fail
//
// ModeTCP:
// a traditional portscan using DialWithTimeout. this performs a full-connect
// scan of the target.
package netscan

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

/* Scanner provides an abstraction for various network scanning methods.
 */
type Scanner interface {
	/* Scan accepts an address in the format of <host>:<port> and returns
	   whether or not that port is open. the method by which it determines
	   this is up to the implementation of the interface.
	*/
	Scan(addr string, timeout time.Duration) bool
}

// Mode is a string alias to refer to the currently supported scanner modes
type Mode string

const (
	/* ModeHTTP refers to a scanning mode that uses HTTP requests/responses
	   to determine if a port is open and serving HTTP content, useful in a
	   situation (like we've run into here) where an IP may report all ports
	   open, but you are interested in locating particular HTTP service.
	*/
	ModeHTTP = "http"

	/* ModeTCP refers to a more traditional mode of port scan that uses tcp
	   connections to determine if a port is open
	*/
	ModeTCP = "tcp"
)

/* ScanJob describes a run of a particular scanner. jobs can consist of multiple
   hosts and ports and will produce a scan of all combinations.
*/
type ScanJob struct {
	scanner Scanner
	mode    Mode
	hosts   []string
	ports   []string
	lock    *semaphore.Weighted
}

// Option is a type alias for a function that takes a ScanJob reference and modifies it
type Option func(*ScanJob)

// WithTLS returns a function that modifies the referenced scanner to enable TLS
func WithTLS() Option {
	return func(s *ScanJob) {
		if s.mode == "http" {
			hs := s.scanner.(*httpScanner)
			hs.tls = true
		}
	}
}

/* New returns a *ScanJob configured with a list of hosts, ports, and any
   additional options. for synchronization, a semaphore is initialized with the
   value 256, which is intended to govern the number of open file handles that a
   scanner can make at a time. this value was chosen by running a ulimit command
   on a macbook, and seemed okay to us...
*/
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

/* Scan executes the defined scan job, iterating over all hosts/port
   combinations and attempting to connect to them using the defined scan mode.
   it returns results over a string channel.
*/
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
