package config

import (
	"regexp"
)

var r, _ = regexp.Compile(`[a-zA-Z]{3}`)

func isValidLocation(location string) bool {
	if len(location) != 3 {
		return false
	}

	matched := r.MatchString(location)

	return matched
}
