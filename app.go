package main

import (
	"log"
	"time"
)

type app struct {
	server   server
	database database
}

func newApp(server server, db database) *app {
	return &app{server, db}
}

func (app *app) run() {
	// Open database connection.
	if err := app.database.Open(); err != nil {
		log.Println("database open error -", err)
	}

	// Run server on a new goroutine.
	go func() {
		if err := app.server.Run(); err != nil {
			log.Println("server run error -", err)
		}
	}()
}

func (app *app) close() {
	var success = true
	if err := app.server.Shutdown(10 * time.Second); err != nil {
		log.Println("Failed to shutdown server -", err)
		success = false
	}
	if err := app.database.Close(); err != nil {
		log.Println("Failed to close database -", err)
		success = false
	}
	if success {
		log.Println("All dependencies closed 👌🏽")
	}
}
