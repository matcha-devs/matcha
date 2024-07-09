// Copyright (c) 2024 Seoyoung Cho.

package internal

import (
	"time"
)

type User struct {
	ID          uint64
	FirstName   string
	MiddleName  string
	LastName    string
	Email       string
	Password    string
	DateOfBirth string
	CreatedOn   time.Time
}

func (user User) IsValid() (valid bool) {
	return user.ID != 0 && "" != user.FirstName && "" != user.LastName && "" != user.Email && "" != user.Password &&
		"" != user.DateOfBirth && user.CreatedOn.Before(time.Now())
}
