package structs

import (
	"errors"
	"strings"
	"unicode"
)

type noun = string

var (
	errEmptyNoun   = errors.New("noun cannot be empty")
	errInvalidNoun = errors.New("noun can only contain letters")
)

func newNoun(input string) (noun, error) {
	if len(input) == 0 {
		return "", errEmptyNoun
	}
	for _, r := range input {
		if !unicode.IsLetter(r) {
			return "", errInvalidNoun
		}
	}
	return input, nil
}

func newProperNoun(input string) (noun, error) {
	return newNoun(strings.ToUpper(input[:1]) + strings.ToLower(input[1:]))
}

func newImproperNoun(input string) (noun, error) {
	return newNoun(strings.ToLower(input))
}
