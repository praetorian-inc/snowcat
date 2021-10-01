package netscan

import (
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type tcpScanner struct{}

func (s *tcpScanner) Scan(target string, timeout time.Duration) bool {
	log.WithFields(log.Fields{
		"addr": target,
	}).Debug("dialing tcp address for port scan")

	conn, err := net.DialTimeout("tcp", target, timeout)

	// i don't like this part...if there's a better way we should do it...
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			return s.Scan(target, timeout)
		}
		return false
	}

	conn.Close()
	return true
}
