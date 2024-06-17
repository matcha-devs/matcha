package main

import "time"

type server interface {
	Run()
	Shutdown(maxClientDisconnectTime time.Duration) error
}

type database interface {
	Open() error
	Close() error
	AuthenticateLogin(username string, password string) error
	AddUser(username string, email string, password string) error
	GetUserID(varName string, variable string) int // TODO(@seoyoungcho213) return errors from all database calls.
	DeleteUser(id int) error
}
