// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/matcha-devs/matcha/internal/mySQL"
	"github.com/matcha-devs/matcha/structs"
)

var (
	deps   *structs.App
	server = http.Server{
		Addr:                         ":8080",
		Handler:                      newMatchaRouter(),
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  time.Second,
		ReadHeaderTimeout:            2 * time.Second,
		WriteTimeout:                 time.Second,
		IdleTimeout:                  30 * time.Second,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}
)

func main() {
	// Create a channel to wait for the "crtl+c" interrupt to ensure dependencies are closed safely before exiting.
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, syscall.SIGINT, syscall.SIGTERM)

	// Open those dependencies now that it is safe to do so.
	deps = structs.NewDeps(mySQL.Open("matcha_db", "internal/mySQL/queries/") )
	defer deps.Close()

	// Run server indefinitely on a new goroutine, then block the main goroutine until ctrl+c interrupt is raised.
	go func() {
		log.Println("Server starting on:", server.Addr, "🫡")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("HTTP server error -", err)
		}
	}()
	<-ctrlC

	// Allow 10s max for clients to disconnect and shut the server down.
	ctx, release := context.WithTimeout(context.Background(), 10*time.Second)
	defer release()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("Server shutdown err -", err)
	}
	log.Println("Server has shutdown 👋🏽")
}
