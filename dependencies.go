package main

import "time"

type server interface {
	Run() error
	Shutdown(maxClientDisconnectTime time.Duration) error
}

type database interface {
	Open() error
	Close() error
	AuthenticateLogin(username string, password string) error
	AddUser(username string, email string, password string) error
	GetUserID(varName string, variable string) (int, error)
	GetUserInfo(id string) (string, string, string, error)
	DeleteUser(id int) error
}
