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

package knownvulns

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestUnmarshallingData(t *testing.T) {
	data := getIstioCVEData()
	assert.NotEqual(t, 0, len(data))
}

func TestGetMatchingVulns(t *testing.T) {
	matchingVulns, _ := GetVulnsForVersion("1.1.0")
	assert.Equal(t, 10, len(matchingVulns))
	matchingVulns, _ = GetVulnsForVersion("1.2.0")
	assert.Equal(t, 11, len(matchingVulns))
	matchingVulns, _ = GetVulnsForVersion("1.3.0")
	assert.Equal(t, 10, len(matchingVulns))
	matchingVulns, _ = GetVulnsForVersion("1.4.0")
	assert.Equal(t, 11, len(matchingVulns))
	matchingVulns, _ = GetVulnsForVersion("1.5.0")
	assert.Equal(t, 10, len(matchingVulns))
	matchingVulns, _ = GetVulnsForVersion("1.8.6")
	assert.Equal(t, 2, len(matchingVulns))
}
