// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	internalDatabase "github.com/matcha-devs/matcha/internal/database"
	internalServer "github.com/matcha-devs/matcha/internal/server"
)

var matcha *app

func main() {
	// Channel to catch "crtl+c" such that dependencies will be closed safely before opening them.
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, syscall.SIGINT, syscall.SIGTERM)

	// err := os.Setenv("MY_SQL_PASSWORD", "Notasqlpassword!")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Open said dependencies.
	matcha = newApp(
		internalServer.New(router()), internalDatabase.New("matcha_db", "root", os.Getenv("MY_SQL_PASSWORD")),
	)
	defer matcha.close()

	// Open database connection.
	if err := matcha.database.Open(); err != nil {
		log.Println("database startup error -", err)
	}

	// Run server on a new goroutine.
	go matcha.server.Run()

	// Block the main goroutine until ctrl+c interrupt is raised.
	<-ctrlC
}
