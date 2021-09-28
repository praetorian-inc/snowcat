package knownvulns

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type VersionRange struct {
	MinVersion uint64
	MaxVersion uint64
}

func (vr *VersionRange) MatchesVersion(version uint64) bool {
	return version <= vr.MaxVersion && version >= vr.MinVersion
}

type IstioCVEData struct {
	AffectedVersions []VersionRange
	DisclosureId     string
	DisclosureUrl    string
	DisclosureDate   string
	ImpactScore      string
	RelatedString    string
}

const BulletinUrl string = "https://istio.io/latest/news/security/"

func convert_string_to_number(versionString string) (uint64, error) {
	versionNumbers := strings.Split(versionString, ".")
	// For now assume we are working with 3 decimals major.minor.revision
	numDecimals := 3
	versionValue := uint64(0)
	for offset, versionNumber := range versionNumbers {
		version, err := strconv.ParseUint(versionNumber, 10, 64)
		if err != nil {
			return 0, err
		}
		multiplier := uint64(math.Pow(10, float64(numDecimals-offset-1)*4))
		versionValue = versionValue + version*multiplier
	}
	return versionValue, nil
}

func parse_affected_versions(affectedVersionsString string) []VersionRange {
	affectedVersionStrings := strings.Split(affectedVersionsString, "<br>")
	versionRanges := []VersionRange{}
	for _, affectedVersionString := range affectedVersionStrings {
		if affectedVersionString == "" {
			continue
		}
		if strings.HasPrefix(affectedVersionString, "All releases prior to ") {
			// Handle All releases prior to 1.9.8
			maxVersion, err := convert_string_to_number(strings.Split(affectedVersionString, "All releases prior to ")[1])
			// we want PRIOR to this version, so drop the revision by 1
			maxVersion = maxVersion - 1
			if err != nil {
				fmt.Println("Found a version string we couldn't convert: " + affectedVersionString)
			}
			vr := VersionRange{
				MinVersion: 0,
				MaxVersion: maxVersion,
			}
			versionRanges = append(versionRanges, vr)
		} else if strings.Contains(affectedVersionString, " to ") {
			// Handle 1.10.0 to 1.10.3
			minVersion, err := convert_string_to_number(strings.Split(affectedVersionString, " to ")[0])
			if err != nil {
				fmt.Println("Found a version string we couldn't convert: " + affectedVersionString)
			}
			maxVersion, err := convert_string_to_number(strings.Split(affectedVersionString, " to ")[1])
			if err != nil {
				fmt.Println("Found a version string we couldn't convert: " + affectedVersionString)
			}
			vr := VersionRange{
				MinVersion: minVersion,
				MaxVersion: maxVersion,
			}
			versionRanges = append(versionRanges, vr)
		} else if strings.HasSuffix(affectedVersionString, "and later") {
			// Ignore CVEs which are applied to ALL versions after a certain point, our auditors check for these issues
			continue

		} else if strings.HasSuffix(affectedVersionString, "patch releases") {
			// Handle edge case of "All 1.8 patch releases"
			minVersion, err := convert_string_to_number(strings.Split(affectedVersionString, " ")[1])
			if err != nil {
				fmt.Println("Found a version string we couldn't convert: " + affectedVersionString)
			}
			maxVersion := minVersion - (minVersion % 10000) + 9999
			vr := VersionRange{
				MinVersion: minVersion,
				MaxVersion: maxVersion,
			}
			versionRanges = append(versionRanges, vr)
		} else {
			// Handle single version number
			version, err := convert_string_to_number(affectedVersionString)
			if err != nil {
				fmt.Println("Found a version string we couldn't convert: " + affectedVersionString)
			}
			vr := VersionRange{
				MinVersion: version,
				MaxVersion: version,
			}
			versionRanges = append(versionRanges, vr)
		}
	}
	return versionRanges
}

func parse_body(body []byte) ([]IstioCVEData, error) {

	cveDataSlice := []IstioCVEData{}

	// First get the table with <table>.*</table>
	r := regexp.MustCompile(`<table>.*</table>`)
	match := r.FindString(string(body))
	if match == "" {
		return nil, fmt.Errorf("Could not find <table>")
	}
	// Next get each row with <tr>.*?</tr>
	r = regexp.MustCompile(`<tr>(.*?)</tr>`)
	matches := r.FindAllStringSubmatch(match, -1)
	for rowNum, stringMatch := range matches {
		// skip the table header
		if rowNum == 0 {
			continue
		}
		// Break it apart by each <td>(.*?)</td>
		tdRegex := regexp.MustCompile(`<td.*?>(.*?)</td>`)

		colMatches := tdRegex.FindAllStringSubmatch(stringMatch[1], -1)

		//<a href=/latest/news/security/istio-security-2019-001/>ISTIO-SECURITY-2019-001</a>
		linkRegex := regexp.MustCompile(`<a href=(.*?)>(.*?)</a>`)
		discMatches := linkRegex.FindAllStringSubmatch(colMatches[0][1], -1)
		impactMatches := linkRegex.FindAllStringSubmatch(colMatches[3][1], -1)
		affectedVersions := parse_affected_versions(colMatches[2][1])
		impactScore := colMatches[3][1]
		if colMatches[3][1] == "" {
			impactScore = "N/A"
		}
		if colMatches[3][1] != "N/A" && impactMatches != nil {
			impactScore = impactMatches[0][2]
		}

		//fmt.Println("Disclosure: " + discMatches[0][2])
		//fmt.Println("Disclosure URL: " + discMatches[0][1])
		//fmt.Println("Date: " + colMatches[1][1])
		//fmt.Println("Affected Releases: " + colMatches[2][1])
		//fmt.Println("Impact Score: " + impactScore)
		//fmt.Println("Related: " + colMatches[4][1])

		data := IstioCVEData{
			AffectedVersions: affectedVersions,
			DisclosureId:     discMatches[0][2],
			DisclosureUrl:    "https://istio.io" + discMatches[0][1],
			DisclosureDate:   colMatches[1][1],
			ImpactScore:      impactScore,
			RelatedString:    colMatches[4][1],
		}
		cveDataSlice = append(cveDataSlice, data)
	}
	return cveDataSlice, nil
}

func scrape_cve_data() ([]IstioCVEData, error) {
	resp, err := http.Get(BulletinUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cveData, err := parse_body(body)
	if err != nil {
		return nil, err
	}
	return cveData, nil
}
