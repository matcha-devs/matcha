// Copyright (c) 2024 Seoyoung Cho.

package internal

import (
	"time"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedOn time.Time
}

func (user User) IsValid() (valid bool) {
	return user.ID != 0 && "" != user.Username && "" != user.Email && "" != user.Password &&
		user.CreatedOn.Before(time.Now())
}
