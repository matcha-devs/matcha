package structs

import (
	"errors"
	"regexp"
	"strconv"
)

type password = string

const minPasswordLen = 12

var (
	errShortPassword    = errors.New("password must have at least " + strconv.Itoa(minPasswordLen) + " characters")
	errSpacesInPassword = errors.New("password must not contain spaces")
	errSimplePassword   = errors.New("password must contain at least 2 of each: lowercase, uppercase, symbols, numbers")
	lowercasePattern    = regexp.MustCompile(`[^\p{Ll}]*\p{Ll}[^\p{Ll}]*\p{Ll}`)
	uppercasePattern    = regexp.MustCompile(`[^\p{Lu}]*\p{Lu}[^\p{Lu}]*\p{Lu}`)
	numberPattern       = regexp.MustCompile(`[^0-9]*[0-9][^0-9]*[0-9]`)
	symbolPattern       = regexp.MustCompile(`[^-\p{L}]*[^\p{L}][^-\p{L}]*[^\p{L}]`)
)

func newPassword(input string) (password, error) {
	if len(input) < minPasswordLen {
		return "", errShortPassword
	}
	if containsSpaces(input) {
		return "", errSpacesInPassword
	}
	if !lowercasePattern.MatchString(input) || !uppercasePattern.MatchString(input) ||
		!numberPattern.MatchString(input) || !symbolPattern.MatchString(input) {
		return "", errSimplePassword
	}

	return input, nil
}
