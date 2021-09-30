/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package namer

import (
	"strings"
)

var consonants = "bcdfghjklmnpqrstvwxyz"

// PluralName returns the plural form of a string.
func PluralName(singular string) string {
	var plural string
	if len(singular) < 2 {
		return strings.ToLower(singular)
	}

	switch rune(singular[len(singular)-1]) {
	case 's', 'x', 'z':
		plural = esPlural(singular)
	case 'y':
		sl := rune(singular[len(singular)-2])
		if isConsonant(sl) {
			plural = iesPlural(singular)
		} else {
			plural = sPlural(singular)
		}
	case 'h':
		sl := rune(singular[len(singular)-2])
		if sl == 'c' || sl == 's' {
			plural = esPlural(singular)
		} else {
			plural = sPlural(singular)
		}
	case 'e':
		sl := rune(singular[len(singular)-2])
		if sl == 'f' {
			plural = vesPlural(singular[:len(singular)-1])
		} else {
			plural = sPlural(singular)
		}
	case 'f':
		plural = vesPlural(singular)
	default:
		plural = sPlural(singular)
	}
	return strings.ToLower(plural)
}

func iesPlural(singular string) string {
	return singular[:len(singular)-1] + "ies"
}

func vesPlural(singular string) string {
	return singular[:len(singular)-1] + "ves"
}

func esPlural(singular string) string {
	return singular + "es"
}

func sPlural(singular string) string {
	return singular + "s"
}

func isConsonant(char rune) bool {
	for _, c := range consonants {
		if char == c {
			return true
		}
	}
	return false
}
