// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/matcha-devs/matcha/internal/sql"
)

var (
	app          = NewApp(sql.Init())
	maxRouteTime = time.Second
	t            = template.Must(
		template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")),
	)
	validEntryPoints = map[string]struct{}{
		"signup": {}, "signup-submit": {}, "signup-fail": {},
		"login": {}, "login-submit": {}, "login-fail": {},
		"dashboard": {}, "settings": {}, "delete-user": {},
	}
)

func loadPage(w http.ResponseWriter, r *http.Request, title string) {
	username := r.FormValue("username")
	user := sql.User{
		ID:       app.db.GetUserID("username", username),
		Username: username,
		Email:    "test",
		Password: "test",
	}
	err := t.ExecuteTemplate(w, title+".gohtml", user)
	if err != nil {
		log.Println("Error executing template - ", err)
	}
}

func signupSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		loadPage(w, r, "/")
		return
	}
	password := r.FormValue("psw")
	if password != r.FormValue("psw-repeat") {
		log.Println("Passwords didnt match.")
		w.WriteHeader(http.StatusBadRequest)
		loadPage(w, r, "signup-fail")
		return
	}
	username := r.FormValue("username")
	err := app.db.AddUser(username, r.FormValue("email"), password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "signup-fail")
	} else {
		loadPage(w, r, "dashboard")
	}
}

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		loadPage(w, r, "/")
		return
	}
	username := r.FormValue("username")
	err := app.db.AuthenticateLogin(username, r.FormValue("password"))
	if err != nil {
		log.Println("Login failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "login-fail")
	} else {
		loadPage(w, r, "dashboard")
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		loadPage(w, r, "/")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	err := app.db.AuthenticateLogin(username, password)
	if err != nil {
		log.Println("Delete User failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "settings")
	} else {
		id := app.db.GetUserID("username", username)
		err := app.db.DeleteUser(id)
		if err != nil {
			log.Println("Delete User failed:", err)
		}
	}
}

func route(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, string(os.PathSeparator))
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
		_, exists := validEntryPoints[path]
		_, err := os.Stat(path)
		switch {
		case exists:
			loadPage(w, r, path)
		case err == nil:
			http.ServeFile(w, r, path)
		case os.IsNotExist(err):
			http.NotFound(w, r)
		default:
			log.Println(err)
		}
	}
}

// TODO(@FaaizMemonPurdue): This is an example of how go routines should be used, but we still need server timeouts
func routeWithTimeout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		log.Println("Routing took longer than", maxRouteTime)
	default:
		start := time.Now()
		route(w, r)
		log.Println("Routing done after", time.Since(start))
	}
}

func main() {
	server := http.Server{
		Addr:                         ":8080",
		Handler:                      http.TimeoutHandler(http.HandlerFunc(routeWithTimeout), maxRouteTime, ""),
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
	log.Println("Starting the server on", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
