package auditors

import (
	"fmt"
	"sync"

	"github.com/praetorian-inc/mithril/pkg/types"
)

func Register(auditor types.Auditor) {
	registryMu.Lock()
	defer registryMu.Unlock()

	name := auditor.Name()
	if _, ok := registry[name]; ok {
		panic(fmt.Errorf("auditor %s already registered", name))
	}
	registry[name] = auditor
}

func New(conf types.Config) ([]types.Auditor, error) {
	registryMu.Lock()
	defer registryMu.Unlock()

	var auditors []types.Auditor
	for _, v := range registry {
		auditors = append(auditors, v)
	}
	return auditors, nil
}

var (
	registry   = make(map[string]types.Auditor)
	registryMu sync.RWMutex
)
