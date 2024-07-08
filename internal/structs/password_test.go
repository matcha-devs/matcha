package structs

import (
	"errors"
	"testing"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		input    string
		expected password
		err      error
	}{
		{"Ab1!Ab1!Ab1!", "Ab1!Ab1!Ab1!", nil},             // Valid password
		{"short1!", "", errShortPassword},                 // Too short
		{"NoNumbersOrSymbolsHere", "", errSimplePassword}, // No numbers or symbols
		{"noHort1!", "", errShortPassword},                // No uppercase and not enough symbols
		{"NoSymbolsHere1", "", errSimplePassword},         // No lowercase or symbols
		{"Ab1!A B1!Ab1!", "", errSpacesInPassword},        // Contains spaces
		{"aB1!aB1!", "", errShortPassword},                // Valid pattern but too short
		{"ab12!@AB", "", errShortPassword},                // Valid pattern but too short
		{"Missing_duplicates1", "", errSimplePassword},    // Missing required elements
		{"aB1!aB1!aB1!aB1!", "aB1!aB1!aB1!aB1!", nil},     // Longer valid password
		{"aB1-aB1_aB12", "aB1-aB1_aB12", nil},             // Valid password with underscores and hyphens
		{"абВГд1!абВГд1!", "абВГд1!абВГд1!", nil},         // Valid password with Unicode letters
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := newPassword(test.input)
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
			if !errors.Is(err, test.err) {
				t.Errorf("expected error %v, got %v", test.err, err)
			}
		})
	}
}
