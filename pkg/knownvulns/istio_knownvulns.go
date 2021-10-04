package knownvulns

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const scrapedCveYamlData = `
- affectedversions:
  - minversion: 0
    maxversion: 100090007
  - minversion: 100100000
    maxversion: 100100003
  - minversion: 100110000
    maxversion: 100110000
  disclosureid: ISTIO-SECURITY-2021-008
  disclosureurl: https://istio.io/latest/news/security/istio-security-2021-008/
  disclosuredate: August 24, 2021
  impactscore: "8.6"
  relatedstring: Multiple CVEs related to AuthorizationPolicy, EnvoyFilter and Envoy
- affectedversions:
  - minversion: 100080000
    maxversion: 100089999
  - minversion: 100090000
    maxversion: 100090005
  - minversion: 100100000
    maxversion: 100100001
  disclosureid: ISTIO-SECURITY-2021-007
  disclosureurl: https://istio.io/latest/news/security/istio-security-2021-007/
  disclosuredate: June 24, 2021
  impactscore: "9.1"
  relatedstring: Istio contains a remotely exploitable vulnerability where credentials
    specified in the Gateway and DestinationRule credentialName field can be accessed
    from different namespaces
- affectedversions:
  - minversion: 0
    maxversion: 100080005
  - minversion: 100090000
    maxversion: 100090004
  disclosureid: ISTIO-SECURITY-2021-005
  disclosureurl: https://istio.io/latest/news/security/istio-security-2021-005/
  disclosuredate: May 11, 2021
  impactscore: "8.1"
  relatedstring: HTTP request paths with multiple slashes or escaped slash characters
    may bypass path based authorization rules
- affectedversions:
  - minversion: 0
    maxversion: 100080005
  - minversion: 100090000
    maxversion: 100090004
  disclosureid: ISTIO-SECURITY-2021-006
  disclosureurl: https://istio.io/latest/news/security/istio-security-2021-006/
  disclosuredate: May 11, 2021
  impactscore: "10"
  relatedstring: An external client can access unexpected services in the cluster,
    bypassing authorization checks, when a gateway is configured with AUTO_PASSTHROUGH
    routing configuration
- affectedversions:
  - minversion: 0
    maxversion: 100080004
  - minversion: 100090000
    maxversion: 100090002
  disclosureid: ISTIO-SECURITY-2021-003
  disclosureurl: https://istio.io/latest/news/security/istio-security-2021-003/
  disclosuredate: April 15, 2021
  impactscore: "7.5"
  relatedstring: ""
- affectedversions: []
  disclosureid: ISTIO-SECURITY-2021-004
  disclosureurl: https://istio.io/latest/news/security/istio-security-2021-004/
  disclosuredate: April 15, 2021
  impactscore: N/A
  relatedstring: Potential misuse of mTLS-only fields in AuthorizationPolicy with
    plain text traffic
- affectedversions: []
  disclosureid: ISTIO-SECURITY-2021-002
  disclosureurl: https://istio.io/latest/news/security/istio-security-2021-002/
  disclosuredate: April 7, 2021
  impactscore: N/A
  relatedstring: Upgrades from older Istio versions can affect access control to an
    ingress gateway due to a change of container ports
- affectedversions:
  - minversion: 100090000
    maxversion: 100090000
  disclosureid: ISTIO-SECURITY-2021-001
  disclosureurl: https://istio.io/latest/news/security/istio-security-2021-001/
  disclosuredate: March 1, 2021
  impactscore: "8.2"
  relatedstring: JWT authentication can be bypassed when AuthorizationPolicy is misused
- affectedversions:
  - minversion: 100080000
    maxversion: 100080000
  disclosureid: ISTIO-SECURITY-2020-011
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-011/
  disclosuredate: November 21, 2020
  impactscore: N/A
  relatedstring: Envoy incorrectly restores the proxy protocol downstream address
    for non-HTTP connections
- affectedversions:
  - minversion: 100060000
    maxversion: 100060010
  - minversion: 100070000
    maxversion: 100070002
  disclosureid: ISTIO-SECURITY-2020-010
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-010/
  disclosuredate: September 29, 2020
  impactscore: "8.3"
  relatedstring: ""
- affectedversions:
  - minversion: 100050000
    maxversion: 100050008
  - minversion: 100060000
    maxversion: 100060007
  disclosureid: ISTIO-SECURITY-2020-009
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-009/
  disclosuredate: August 11, 2020
  impactscore: "6.8"
  relatedstring: Incorrect Envoy configuration for wildcard suffixes used for Principals/Namespaces
    in Authorization Policies for TCP Services
- affectedversions:
  - minversion: 100050000
    maxversion: 100050007
  - minversion: 100060000
    maxversion: 100060004
  - minversion: 0
    maxversion: 100049999
  disclosureid: ISTIO-SECURITY-2020-008
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-008/
  disclosuredate: July 9, 2020
  impactscore: "6.6"
  relatedstring: Incorrect validation of wildcard DNS Subject Alternative Names
- affectedversions:
  - minversion: 100050000
    maxversion: 100050006
  - minversion: 100060000
    maxversion: 100060003
  disclosureid: ISTIO-SECURITY-2020-007
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-007/
  disclosuredate: June 30, 2020
  impactscore: "7.5"
  relatedstring: Multiple denial of service vulnerabilities in Envoy
- affectedversions:
  - minversion: 100040000
    maxversion: 100040009
  - minversion: 100050000
    maxversion: 100050004
  - minversion: 100060000
    maxversion: 100060001
  disclosureid: ISTIO-SECURITY-2020-006
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-006/
  disclosuredate: June 11, 2020
  impactscore: "7.5"
  relatedstring: Denial of service in the HTTP2 library used by Envoy
- affectedversions:
  - minversion: 100040000
    maxversion: 100040008
  - minversion: 100050000
    maxversion: 100050003
  disclosureid: ISTIO-SECURITY-2020-005
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-005/
  disclosuredate: May 12, 2020
  impactscore: "7.5"
  relatedstring: Denial of service affecting telemetry v2
- affectedversions:
  - minversion: 100040000
    maxversion: 100040006
  - minversion: 100050000
    maxversion: 100050000
  disclosureid: ISTIO-SECURITY-2020-004
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-004/
  disclosuredate: March 25, 2020
  impactscore: "8.7"
  relatedstring: Default Kiali security configuration allows full control of mesh
- affectedversions:
  - minversion: 100040000
    maxversion: 100040005
  disclosureid: ISTIO-SECURITY-2020-003
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-003/
  disclosuredate: March 3, 2020
  impactscore: "7.5"
  relatedstring: Two uncontrolled resource consumption and two incorrect access control
    vulnerabilities in Envoy
- affectedversions:
  - minversion: 100030000
    maxversion: 100030006
  disclosureid: ISTIO-SECURITY-2020-002
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-002/
  disclosuredate: February 11, 2020
  impactscore: "7.4"
  relatedstring: Mixer policy check bypass caused by improperly accepting certain
    request headers
- affectedversions:
  - minversion: 100030000
    maxversion: 100030007
  - minversion: 100040000
    maxversion: 100040003
  disclosureid: ISTIO-SECURITY-2020-001
  disclosureurl: https://istio.io/latest/news/security/istio-security-2020-001/
  disclosuredate: February 11, 2020
  impactscore: "9.0"
  relatedstring: Authentication Policy bypass
- affectedversions:
  - minversion: 100020000
    maxversion: 100020009
  - minversion: 100030000
    maxversion: 100030005
  - minversion: 100040000
    maxversion: 100040001
  disclosureid: ISTIO-SECURITY-2019-007
  disclosureurl: https://istio.io/latest/news/security/istio-security-2019-007/
  disclosuredate: December 10, 2019
  impactscore: "9.0"
  relatedstring: Heap overflow and improper input validation in Envoy
- affectedversions:
  - minversion: 100030000
    maxversion: 100030004
  disclosureid: ISTIO-SECURITY-2019-006
  disclosureurl: https://istio.io/latest/news/security/istio-security-2019-006/
  disclosuredate: November 7, 2019
  impactscore: "7.5"
  relatedstring: Denial of service
- affectedversions:
  - minversion: 100010000
    maxversion: 100010015
  - minversion: 100020000
    maxversion: 100020006
  - minversion: 100030000
    maxversion: 100030001
  disclosureid: ISTIO-SECURITY-2019-005
  disclosureurl: https://istio.io/latest/news/security/istio-security-2019-005/
  disclosuredate: October 8, 2019
  impactscore: "7.5"
  relatedstring: Denial of service caused by the presence of numerous HTTP headers
    in client requests
- affectedversions:
  - minversion: 100020000
    maxversion: 100020004
  disclosureid: Istio 1.2.4 sidecar image vulnerability
  disclosureurl: https://istio.io/latest/news/security/incorrect-sidecar-image-1.2.4/
  disclosuredate: September 10, 2019
  impactscore: N/A
  relatedstring: An erroneous 1.2.4 sidecar image was available due to a faulty release
    operation
- affectedversions:
  - minversion: 100010000
    maxversion: 100010012
  - minversion: 100020000
    maxversion: 100020003
  disclosureid: ISTIO-SECURITY-2019-003
  disclosureurl: https://istio.io/latest/news/security/istio-security-2019-003/
  disclosuredate: August 13, 2019
  impactscore: "7.5"
  relatedstring: Denial of service in regular expression parsing
- affectedversions:
  - minversion: 100010000
    maxversion: 100010012
  - minversion: 100020000
    maxversion: 100020003
  disclosureid: ISTIO-SECURITY-2019-004
  disclosureurl: https://istio.io/latest/news/security/istio-security-2019-004/
  disclosuredate: August 13, 2019
  impactscore: "7.5"
  relatedstring: Multiple denial of service vulnerabilities related to HTTP2 support
    in Envoy
- affectedversions:
  - minversion: 100000000
    maxversion: 100000008
  - minversion: 100010000
    maxversion: 100010009
  - minversion: 100020000
    maxversion: 100020001
  disclosureid: ISTIO-SECURITY-2019-002
  disclosureurl: https://istio.io/latest/news/security/istio-security-2019-002/
  disclosuredate: June 28, 2019
  impactscore: "7.5"
  relatedstring: Denial of service affecting JWT access token parsing
- affectedversions:
  - minversion: 100010000
    maxversion: 100010006
  disclosureid: ISTIO-SECURITY-2019-001
  disclosureurl: https://istio.io/latest/news/security/istio-security-2019-001/
  disclosuredate: May 28, 2019
  impactscore: "8.9"
  relatedstring: Incorrect access control`

func getIstioCVEData() []IstioCVEData {
	var data []IstioCVEData
	err := yaml.Unmarshal([]byte(scrapedCveYamlData), &data)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return data
}

// GetVulnsForVersion returns an array of Istio CVEs matching a given version.
func GetVulnsForVersion(version string) ([]IstioCVEData, error) {
	vulnData := getIstioCVEData()
	results := []IstioCVEData{}
	versionNum, err := convertStringToNumber(version)
	if err != nil {
		return nil, err
	}
	for _, vuln := range vulnData {
		for _, versionRange := range vuln.AffectedVersions {
			if versionRange.MatchesVersion(versionNum) {
				results = append(results, vuln)
				break
			}
		}
	}
	return results, nil
}
