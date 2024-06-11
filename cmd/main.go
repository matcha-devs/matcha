// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/matcha-devs/matcha/internal/database"

	_ "github.com/go-sql-driver/mysql"
)

var (
	validEntryPoints = map[string]struct{}{
		"signup": {}, "signup-submit": {}, "signup-fail": {},
		"login": {}, "login-submit": {}, "login-fail": {},
		"dashboard": {}, "settings": {}, "delete-user": {},
	}
	t = template.Must(
		template.ParseGlob(strings.Join([]string{"internal", "templates", "*.html"}, string(os.PathSeparator))),
	)
	maxRouteTime = time.Second
)

func loadPage(w http.ResponseWriter, r *http.Request, title string) {
	username := r.FormValue("username")
	user := database.User{
		ID:       database.GetUserID("username", username),
		Username: username,
		Email:    "test",
		Password: "test",
	}
	err := t.ExecuteTemplate(w, title+".html", user)
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
		fmt.Println("Passwords didnt match.")
		w.WriteHeader(http.StatusBadRequest)
		loadPage(w, r, "signup-fail")
		return
	}
	username := r.FormValue("username")
	err := database.AddUser(username, r.FormValue("email"), password)
	fmt.Print()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "signup-fail")
	} else {
		//http.Redirect(w, r, "dashboard?username="+username, http.StatusFound)
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
	err := database.AuthenticateLogin(username, r.FormValue("password"))
	if err != nil {
		log.Println("Login failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "login-fail")
	} else {
		//http.Redirect(w, r, "dashboard?username="+username, http.StatusFound)
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
	err := database.AuthenticateLogin(username, password)
	if err != nil {
		log.Println("Delete User failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "settings")
	} else {
		id := database.GetUserID("username", username)
		database.DeleteUser(id)
	}
}

func route(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimLeft(r.URL.Path, "/")
	fmt.Println(title)
	switch title {
	case "signup-submit":
		signupSubmit(w, r)
	case "login-submit":
		loginSubmit(w, r)
	case "delete-user":
		deleteUser(w, r)
	case "":
		loadPage(w, r, "landing")
	case "main.css":
		http.ServeFile(w, r, "internal/templates/main.css")
	default:
		if _, exists := validEntryPoints[title]; exists {
			loadPage(w, r, title)
		} else {
			http.NotFound(w, r)
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
	database.Init()
	server := http.Server{
		Addr:         ":8080",
		WriteTimeout: 5 * time.Second,
		Handler:      http.TimeoutHandler(http.HandlerFunc(routeWithTimeout), maxRouteTime, "Timeout!\n"),
	}
	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
