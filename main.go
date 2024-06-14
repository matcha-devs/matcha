// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/matcha-devs/matcha/internal/mySQL"
)

var (
	deps         = NewDeps(mySQL.Open())
	maxRouteTime = time.Second
	tmpl         = template.Must(
		template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")),
	)
	validEntryPoints = map[string]struct{}{
		"signup": {}, "signup-submit": {}, "signup-fail": {},
		"login": {}, "login-submit": {}, "login-fail": {},
		"dashboard": {}, "settings": {}, "delete-user": {},
	}
)

func route(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	log.Println("Routing {" + path + "}")
	switch path {
	case "":
		loadPage(w, r, "index")
	case "signup-submit":
		signupSubmit(w, r)
	case "login-submit":
		loginSubmit(w, r)
	case "delete-user":
		deleteUser(w, r)
	default:
		if _, exists := validEntryPoints[path]; exists {
			loadPage(w, r, path)
		} else {
			http.NotFound(w, r)
		}
	}
}

// TODO(@FaaizMemonPurdue): This is an example of how go routines should be used, but we still need server timeouts
func routeWithTimeout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), maxRouteTime)
	defer cancel()
	select {
	case <-ctx.Done():
		log.Println("Routing took longer than", maxRouteTime)
	//case <-time.After(maxRouteTime):
	default:
		start := time.Now()
		route(w, r)
		log.Println("Routing done after", time.Since(start))
	}
}

func main() {
	// Close dependencies before exiting main
	defer deps.Close()

	// Create a waiting channel to catch the "crtl+c" interrupt
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	// TODO(@CarlosACJ55): Make a clean transition from the switch case to ServeMux
	//mux.Handle("/{$}", http.TimeoutHandler(http.HandlerFunc(loadPage), maxRouteTime, ""))
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	mux.Handle("/", http.TimeoutHandler(http.HandlerFunc(routeWithTimeout), maxRouteTime, ""))
	server := http.Server{
		Addr:                         ":8080",
		Handler:                      mux,
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
	go func() {
		log.Println("Starting the server on", server.Addr, "🫡")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("HTTP server error", err)
		}
	}()
	<-ctrlC
	ctx, release := context.WithTimeout(context.Background(), 10*time.Second)
	defer release()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("Server shutdown err -", err)
	}
	log.Println("Server shutdown 😴")
}
