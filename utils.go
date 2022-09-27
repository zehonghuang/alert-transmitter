package main

import "regexp"

func IsBlank(s string) bool {
	if s == "" {
		return true
	}
	if regexp.MustCompile(`^\s+$`).MatchString(s) {
		return true
	}
	return false
}
