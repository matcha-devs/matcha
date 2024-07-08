package structs

import "testing"

func equalName(a, b name) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNewName(t *testing.T) {
	tests := []struct {
		input    string
		expected name
		err      error
	}{
		{input: "John Doe", expected: name{"John", "Doe"}, err: nil},
		{input: "Alice Wonderland", expected: name{"Alice", "Wonderland"}, err: nil},
		{input: "  LeadingSpace", expected: nil, err: errInvalidName},
		{input: "TrailingSpace  ", expected: nil, err: errInvalidName},
		{input: "Double  Space", expected: nil, err: errInvalidName},
		{input: "", expected: nil, err: errEmptyName},
		{input: "John123 Doe", expected: nil, err: errInvalidNoun},
	}

	for _, test := range tests {
		result, err := newName(test.input)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("newName(%q) returned error %v, expected %v", test.input, err, test.err)
		} else if err == nil && test.err != nil {
			t.Errorf("newName(%q) did not return error, expected %v", test.input, test.err)
		} else if !equalName(result, test.expected) {
			t.Errorf("newName(%q) returned %v, expected %v", test.input, result, test.expected)
		}
	}
}
