package main

import "log"

type Deps struct {
	DB Database
}

func NewDeps(db Database) *Deps {
	return &Deps{db}
}

func (deps Deps) Close() {
	err := deps.DB.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("Dependencies closed 🩺")
}

type Database interface {
	Close() error
	AuthenticateLogin(username string, password string) error
	AddUser(username string, email string, password string) error
	GetUserID(varName string, variable string) int
	DeleteUser(id int) error
}
