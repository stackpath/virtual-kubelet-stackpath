package config

import (
	"regexp"

	uuid "github.com/satori/go.uuid"
)

var r, _ = regexp.Compile(`[a-zA-Z]{3}`)

func isValidUUID(uuidStr string) bool {
	_, err := uuid.FromString(uuidStr)
	return err == nil
}

func isValidLocation(location string) bool {
	if len(location) != 3 {
		return false
	}

	matched := r.MatchString(location)

	return matched
}
