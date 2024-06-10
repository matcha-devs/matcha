// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CarlosACJ55/matcha/internal/database"

	_ "github.com/go-sql-driver/mysql"
)

var (
	validEntryPoints = map[string]struct{}{
		"signup": {}, "signup-submit": {}, "signup-fail": {},
		"login": {}, "login-submit": {}, "login-fail": {},
		"dashboard": {},
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
	time.Sleep(3 * time.Second)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	password := r.FormValue("psw")
	if password != r.FormValue("psw-repeat") {
		fmt.Println("Passwords didnt match.")
		http.Redirect(w, r, "signup-fail", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	err := database.AddUser(username, r.FormValue("email"), password)
	if err != nil {
		http.Redirect(w, r, "signup-fail", http.StatusUnauthorized)
	} else {
		http.Redirect(w, r, "/dashboard?username="+username, http.StatusFound)
	}
}

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	err := database.AuthenticateLogin(username, r.FormValue("password"))
	if err != nil {
		log.Println("Login failed:", err)
		http.Redirect(w, r, "login-fail", http.StatusUnauthorized) //TODO(@andreag0101): fix failure redirects
	} else {
		http.Redirect(w, r, "dashboard?username="+username, http.StatusFound)
	}
}

func route(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimLeft(r.URL.Path, "/")
	switch title {
	case "signup-submit":
		signupSubmit(w, r)
	case "login-submit":
		loginSubmit(w, r)
	case "delete-user":
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			panic(err)
		}
		database.DeleteUser(id)
	case "":
		loadPage(w, r, "landing")
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
