package main

import (
	"time"

	"github.com/matcha-devs/matcha/internal"
)

type server interface {
	Run() error
	Shutdown(maxClientDisconnectTime time.Duration) error
}

type database interface {
	Open() error
	Close() error
	AuthenticateLogin(username string, password string) (id int, err error)
	GetUser(id int) (user *internal.User)
	AddUser(username string, email string, password string) error
	GetUserID(varName string, variable string) (id int)
	DeleteUser(id int) error
}
