package util

import "strings"

func ContainsAny(testString string, checkStrings ...string) bool {
	for _, str := range checkStrings {
		if strings.Contains(testString, str) {
			return true
		}
	}
	return false
}
