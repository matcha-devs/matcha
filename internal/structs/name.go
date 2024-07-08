package structs

import (
	"errors"
	"strings"
)

type name []noun

var (
	errEmptyName   = errors.New("name cannot be empty")
	errInvalidName = errors.New("name must not lead or trail with space, nor have double spaces")
)

func newName(input string) (name name, err error) {
	if len(input) == 0 {
		return nil, errEmptyName
	}
	words := strings.Split(input, " ")
	for _, word := range words {
		if len(word) == 0 {
			return nil, errInvalidName
		}
		n, err := newProperNoun(word)
		if err != nil {
			return nil, err
		}
		name = append(name, n)
	}
	return
}

func (n name) String() string {
	return strings.Join(n, " ")
}
