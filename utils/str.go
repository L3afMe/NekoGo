package utils

import "strings"

func SAContainsCI(h []string, j string) bool {
	for _, k := range h {
		if strings.EqualFold(j, k) {
			return true
		}
	}
	return false
}

func SAContains(h []string, j string) bool {
	for _, k := range h {
		if j == k {
			return true
		}
	}
	return false
}

func SContainsCI(h string, j string) bool {
	return strings.Contains(strings.ToLower(h), strings.ToLower(j))
}
