package structs

import (
	"errors"
	"testing"
	"time"
)

func TestNewDateOfBirth(t *testing.T) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	maxDateAgo := today.AddDate(-maxUserAge, 0, 0)
	dayBeforeMaxDateAgo := maxDateAgo.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)
	tests := []struct {
		name     string
		input    string
		expected dateOfBirth
		err      error
	}{
		{"valid_today", today.Format(time.DateOnly), today, nil},
		{"valid_200_years_ago", maxDateAgo.Format(time.DateOnly), maxDateAgo, nil},
		{"invalid_201_years_ago", dayBeforeMaxDateAgo.Format(time.DateOnly), dateOfBirth{}, errInvalidDateOfBirth},
		{"invalid_tomorrow", tomorrow.Format(time.DateOnly), dateOfBirth{}, errInvalidDateOfBirth},
		{"invalid_feb_day", "2023-02-30", dateOfBirth{}, errMalformedDateOfBirth},
		{"invalid_month", "2023-13-01", dateOfBirth{}, errMalformedDateOfBirth},
		{"invalid_format", "01-01-2023", dateOfBirth{}, errMalformedDateOfBirth},
		{"invalid_string", "2023-01-01 extra", dateOfBirth{}, errMalformedDateOfBirth},
		{"invalid_extra_info", "2023-01-01T12:00:00Z", dateOfBirth{}, errMalformedDateOfBirth},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := newDateOfBirth(test.input)
			if !result.Equal(test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
			if !errors.Is(err, test.err) {
				t.Errorf("expected error %v, got %v", test.err, err)
			}
		})
	}
}
