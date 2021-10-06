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

// Package auditors defines an interface for all auditors.
// This package provides a registration mechanism similar to
// database/sql. In order to use a particular auditor, add a
// "blank import" for its package.
//
//  import (
//      "github.com/praetorian-inc/mithril/auditors"
//      "github.com/praetorian-inc/mithril/pkg/types"
//
//      _ "github.com/praetorian-inc/mithril/auditors/authz"
//      _ "github.com/praetorian-inc/mithril/auditors/peerauth"
//  )
//
//  auditors, err := auditors.New(types.Config{})
//  if err != nil {
//      // handle error
//  }
//  for _, auditor := range auditors {
//      res, err := auditor.Audit(disco, resources)
//      ...
//  }
//
// See https://golang.org/doc/effective_go.html#blank_import for more
// information on "blank imports".
package auditors

import (
	"fmt"
	"sync"

	"github.com/praetorian-inc/mithril/pkg/types"
)

// Register makes an auditor available with the provided name. If register is
// called twice or if the driver is nil, if panics. Register() is typically
// called in the auditor implementation's init() function to allow for easy
// importing of each auditor.
func Register(auditor types.Auditor) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if auditor == nil {
		panic("Registered auditor is nil")
	}

	name := auditor.Name()
	if _, ok := registry[name]; ok {
		panic(fmt.Errorf("auditor %s already registered", name))
	}
	registry[name] = auditor
}

// All returns a list of all auditors.
func All() []types.Auditor {
	registryMu.Lock()
	defer registryMu.Unlock()

	var auditors []types.Auditor
	for _, v := range registry {
		auditors = append(auditors, v)
	}
	return auditors
}

var (
	registry   = make(map[string]types.Auditor)
	registryMu sync.RWMutex
)
