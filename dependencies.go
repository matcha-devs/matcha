package main

type Dependencies struct {
	db Database
}

func InitDependencies(db Database) *Dependencies {
	return &Dependencies{db}
}

type Database interface {
	Close() error
	AuthenticateLogin(username string, password string) error
	AddUser(username string, email string, password string) error
	GetUserID(varName string, variable string) int
	DeleteUser(id int) error
}
