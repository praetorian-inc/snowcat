package knownvulns

import (
	"testing"

	"github.com/bmizerany/assert"
	"gopkg.in/yaml.v2"
)

func TestVersionConversion(t *testing.T) {
	version, _ := convertStringToNumber("1.2.3")
	assert.Equal(t, uint64(100020003), version)
	version, _ = convertStringToNumber("1234.5678.9012")
	assert.Equal(t, uint64(123456789012), version)

	minVersion, _ := convertStringToNumber("1.1.8")
	maxVersion, _ := convertStringToNumber("1.1.12")
	targetVersion, _ := convertStringToNumber("1.1.10")
	assert.Equal(t, true, minVersion <= targetVersion)
	assert.Equal(t, true, targetVersion <= maxVersion)

	vr := VersionRange{
		MinVersion: minVersion,
		MaxVersion: maxVersion,
	}

	assert.Equal(t, true, vr.MatchesVersion(targetVersion))
	assert.Equal(t, false, vr.MatchesVersion(0))
}

func TestCveScraper(t *testing.T) {
	scrapedData, err := scrapeCVEs()
	if err == nil {
		yamlData, err := yaml.Marshal(scrapedData)
		if err != nil {
			t.Fatalf("%v", err)
		}
		sYamlData := string(yamlData)
		t.Logf("%s", sYamlData)
	}
}
