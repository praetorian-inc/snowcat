package knownvulns

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestUnmarshallingData(t *testing.T) {
	data := GetIstioCVEData()
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
