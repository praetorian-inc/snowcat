package netscan

import (
	"net"
	"strings"
	"time"
)

type tcpScanner struct{}

func (s *tcpScanner) Scan(target string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", target, timeout)

	// i don't like this part...if there's a better way we should do it...
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			return s.Scan(target, timeout)
		} else {
			return false
		}
	}

	conn.Close()
	return true
}
