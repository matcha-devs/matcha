// Copyright (c) 2024 Seoyoung Cho.

package structs

import (
	"errors"
	"time"
)

var ErrMalformedMiddleName = errors.New("unexpected empty middle name")

type CreateUserInput struct {
	isValid       bool
	firstName     name
	hasMiddleName bool
	middleName    name
	lastName      name
	email         email
	password      password
	dateOfBirth   dateOfBirth
}

func NewCreateUserInput(
	hasMiddleName bool, firstName, middleName, lastName, email, password, dateOfBirth string,
) (new *CreateUserInput, err error) {
	new.firstName, err = newName(firstName)
	if new.middleName, err = newName(middleName); errors.Is(err, errEmptyName) == hasMiddleName {
		err = ErrMalformedMiddleName
	}
	new.lastName, err = newName(lastName)
	new.email, err = newEmail(email)
	new.password, err = newPassword(password)
	new.dateOfBirth, err = newDateOfBirth(dateOfBirth)
	new.isValid = err == nil
	return
}

func (input *CreateUserInput) FirstName() string {
	return input.firstName.String()
}

func (input *CreateUserInput) MiddleName() (string, bool) {
	return input.middleName.String(), input.hasMiddleName
}

func (input *CreateUserInput) LastName() string {
	return input.lastName.String()
}

func (input *CreateUserInput) Email() string {
	return input.email
}

func (input *CreateUserInput) Password() string {
	return input.password
}

func (input *CreateUserInput) DateOfBirth() time.Time {
	return input.dateOfBirth
}
