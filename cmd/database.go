package main

type databaseContract interface {
	Close() error
	AuthenticateLogin(username string, password string) error
	AddUser(username string, email string, password string) error
	GetUserID(varName string, variable string) int
	DeleteUser(id int) error
}

type Application struct {
	db databaseContract
}

func NewApplication(db databaseContract) *Application {
	return &Application{db}
}
