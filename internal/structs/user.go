// Copyright (c) 2024 Seoyoung Cho.

package structs

import "time"

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}
