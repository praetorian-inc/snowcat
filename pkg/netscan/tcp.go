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
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type tcpScanner struct{}

func (s *tcpScanner) Scan(target string, timeout time.Duration) bool {
	log.WithFields(log.Fields{
		"addr": target,
	}).Trace("dialing tcp address for port scan")

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
