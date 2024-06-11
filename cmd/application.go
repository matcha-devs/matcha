package main

type Application struct {
	db DatabaseInterface
}

func NewApplication(db DatabaseInterface) *Application {
	return &Application{db}
}
