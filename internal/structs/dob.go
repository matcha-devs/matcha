package structs

import (
	"errors"
	"strconv"
	"time"
)

type dateOfBirth = time.Time

const (
	maxUserAge        = 200
	dateOfBirthLayout = "2006-01-02"
)

var (
	errMalformedDateOfBirth = errors.New("date of birth must be in the form:" + dateOfBirthLayout)
	errInvalidDateOfBirth   = errors.New("date of birth must be less than " + strconv.Itoa(maxUserAge) + " years ago")
)

func newDateOfBirth(input string) (dateOfBirth, error) {
	date, err := time.Parse(dateOfBirthLayout, input)
	if err != nil {
		return dateOfBirth{}, errMalformedDateOfBirth
	}
	if date.After(time.Now().UTC()) || date.Before(time.Now().UTC().AddDate(-maxUserAge, 0, -1)) {
		return dateOfBirth{}, errInvalidDateOfBirth
	}
	return date, nil
}
