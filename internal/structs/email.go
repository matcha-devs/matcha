package structs

import (
	"errors"
	"strings"
	"unicode"
)

type email = string

var (
	errSpacesInEmail = errors.New("email must not contain spaces")
	errInvalidEmail  = errors.New("email must have exactly 1 @")
)

func containsSpaces(input string) bool {
	return strings.IndexFunc(input, unicode.IsSpace) != -1
}

func newEmail(input string) (email, error) {
	words := strings.Split(input, "@")
	if containsSpaces(input) {
		return "", errSpacesInEmail
	}
	if len(words) != 2 {
		return "", errInvalidEmail
	}
	return input, nil
}
