// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"os"
	"os/signal"
	"syscall"

	internalDatabase "github.com/matcha-devs/matcha/internal/database"
	internalServer "github.com/matcha-devs/matcha/internal/server"
)

var matcha *app

func init() {
	matcha = newApp(
		internalServer.New(loggedRouter()), internalDatabase.New("matcha_db", "root", os.Getenv("MYSQL_PASSWORD")),
	)
}

func main() {
	// Channel to catch "crtl+c" such that dependencies will be closed safely before opening them.
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, syscall.SIGINT, syscall.SIGTERM)

	// Open said dependencies and run application on a new go routine.
	go matcha.run()

	// Block the main goroutine until ctrl+c interrupt is raised.
	<-ctrlC

	// Stop application and close the dependencies before exiting.
	matcha.close()
}
