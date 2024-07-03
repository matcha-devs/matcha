// Copyright (c) 2024 Seoyoung Cho.

package internal

import (
	"time"
)

type User struct {
	ID         int
	Firstname  string
	Middlename string
	Lastname   string
	Email      string
	Password   string
	Birthdate  string
	CreatedOn  time.Time
}

func (user User) IsValid() (valid bool) {
	return user.ID != 0 && "" != user.Firstname && "" != user.Lastname && "" != user.Email && "" != user.Password &&
		"" != user.Birthdate && user.CreatedOn.Before(time.Now())
}
