package utils

import "strings"

func Contains(elems []string, v string) bool {
	for _, s := range elems {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}
