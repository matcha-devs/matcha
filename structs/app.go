package structs

import (
	"log"

	"github.com/matcha-devs/matcha/dependencies"
)

type App struct {
	DB dependencies.Database
}

func NewDeps(db dependencies.Database) *App {
	return &App{db}
}

func (deps App) Close() {
	err := deps.DB.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("Dependencies closed 👌🏽")
}
