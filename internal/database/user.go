// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package database

import "time"

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}
