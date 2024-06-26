// Copyright (c) 2024 Seoyoung Cho.

package internal

import (
	"database/sql"
)

type User struct {
	ID        sql.Null[int]
	Username  sql.NullString
	Email     sql.NullString
	Password  sql.NullString
	CreatedOn sql.NullTime
}

func (user User) IsValid() (valid bool) {
	return user.ID.Valid && user.Username.Valid && user.Email.Valid && user.Password.Valid && user.CreatedOn.Valid
}
