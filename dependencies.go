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
	AuthenticateLogin(email string, password string) (id int, err error)
	GetUser(id int) (user *internal.User)
	AddUser(firstname string, middlename string, lastname string, email string, password string,
		birthdate string) (id int, err error)
	GetUserID(varName string, variable string) (id int)
	DeleteUser(id int) (err error)
}
