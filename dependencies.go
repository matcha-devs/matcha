package main

import (
	"time"

	"github.com/matcha-devs/matcha/internal"
)

type server interface {
	Run() (err error)
	Shutdown(maxClientDisconnectTime time.Duration) (err error)
}

type database interface {
	Open() (err error)
	Close() (err error)
	AuthenticateLogin(email, password string) (id uint64, err error)
	GetUser(id uint64) (user *internal.User)
	AddUser(first_name, middle_name, last_name, email, password, date_of_birth string) (id uint64, err error)
	GetUserID(email string) (id uint64)
	DeleteUser(id uint64) (err error)
}
