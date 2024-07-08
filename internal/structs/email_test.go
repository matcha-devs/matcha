package structs

import "testing"

func TestNewEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected email
		err      error
	}{
		{input: "user@example.com", expected: email("user@example.com"), err: nil},
		{input: "user.name@example.com", expected: email("user.name@example.com"), err: nil},
		{input: "user@example", expected: email("user@example"), err: nil},
		{input: "user@.com", expected: email("user@.com"), err: nil},
		{input: "user @example.com", expected: "", err: errSpacesInEmail},
		{input: "user@example.com ", expected: "", err: errSpacesInEmail},
		{input: "user@subdomain.example.com", expected: email("user@subdomain.example.com"), err: nil},
		{input: "user@local@host", expected: "", err: errInvalidEmail},
		{input: "user@subdomain.example.com", expected: email("user@subdomain.example.com"), err: nil},
		{input: "", expected: "", err: errInvalidEmail},
	}

	for _, test := range tests {
		result, err := newEmail(test.input)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("newEmail(%q) returned error %v, expected %v", test.input, err, test.err)
		} else if err == nil && test.err != nil {
			t.Errorf("newEmail(%q) did not return error, expected %v", test.input, test.err)
		} else if result != test.expected {
			t.Errorf("newEmail(%q) returned %v, expected %v", test.input, result, test.expected)
		}
	}
}
