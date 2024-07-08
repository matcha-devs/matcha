package structs

import (
	"testing"
)

func TestNewNoun(t *testing.T) {
	tests := []struct {
		input    string
		expected noun
		err      error
	}{
		{input: "Alice", expected: "Alice", err: nil},
		{input: "Bob", expected: "Bob", err: nil},
		{input: "", expected: "", err: errEmptyNoun},
		{input: "John123", expected: "", err: errInvalidNoun},
		{input: "Mary ", expected: "", err: errInvalidNoun},
		{input: "  Mike", expected: "", err: errInvalidNoun},
		{input: "Anne-Marie", expected: "", err: errInvalidNoun},
	}

	for _, test := range tests {
		result, err := newNoun(test.input)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("newNoun(%q) returned error %v, expected %v", test.input, err, test.err)
		} else if err == nil && test.err != nil {
			t.Errorf("newNoun(%q) did not return error, expected %v", test.input, test.err)
		} else if result != test.expected {
			t.Errorf("newNoun(%q) returned %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestNewProperNoun(t *testing.T) {
	tests := []struct {
		input    string
		expected noun
		err      error
	}{
		{input: "alice", expected: "Alice", err: nil},
		{input: "bob", expected: "Bob", err: nil},
		{input: "  leading", expected: "", err: errInvalidNoun},
		{input: "trailing  ", expected: "", err: errInvalidNoun},
		{input: "double  space", expected: "", err: errInvalidNoun},
		{input: "John123", expected: "", err: errInvalidNoun},
		{input: "   spaces", expected: "", err: errInvalidNoun},
	}

	for _, test := range tests {
		result, err := newProperNoun(test.input)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("newProperNoun(%q) returned error %v, expected %v", test.input, err, test.err)
		} else if err == nil && test.err != nil {
			t.Errorf("newProperNoun(%q) did not return error, expected %v", test.input, test.err)
		} else if result != test.expected {
			t.Errorf("newProperNoun(%q) returned %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestNewImproperNoun(t *testing.T) {
	tests := []struct {
		input    string
		expected noun
		err      error
	}{
		{input: "ALICE", expected: "alice", err: nil},
		{input: "BOB", expected: "bob", err: nil},
		{input: "  leading", expected: "", err: errInvalidNoun},
		{input: "trailing  ", expected: "", err: errInvalidNoun},
		{input: "double  space", expected: "", err: errInvalidNoun},
		{input: "John123", expected: "", err: errInvalidNoun},
		{input: "   spaces", expected: "", err: errInvalidNoun},
	}

	for _, test := range tests {
		result, err := newImproperNoun(test.input)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("newImproperNoun(%q) returned error %v, expected %v", test.input, err, test.err)
		} else if err == nil && test.err != nil {
			t.Errorf("newImproperNoun(%q) did not return error, expected %v", test.input, test.err)
		} else if result != test.expected {
			t.Errorf("newImproperNoun(%q) returned %v, expected %v", test.input, result, test.expected)
		}
	}
}
