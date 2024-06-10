// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/CarlosACJ55/matcha/internal/database"

	_ "github.com/go-sql-driver/mysql"
)

var (
	validEntryPoints = map[string]struct{}{
		"login": {}, "signup": {}, "dashboard": {}, "login_fail": {},
		"signup_fail": {},
	}
	t = template.Must(
		template.ParseGlob(strings.Join([]string{"internal", "templates", "*.html"}, string(os.PathSeparator))),
	)
)

func signupSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	password := r.FormValue("psw")
	if password != r.FormValue("psw_repeat") {
		http.Redirect(w, r, "signup_fail", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	err := database.AddUser(username, r.FormValue("email"), password)
	if err != nil {
		http.Redirect(w, r, "login_fail", http.StatusUnauthorized)
	} else {
		http.Redirect(w, r, "/dashboard?username="+username, http.StatusFound)
	}
}

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	err := database.AuthenticateLogin(username, r.FormValue("password"))
	if err != nil {
		http.Redirect(w, r, "/login_fail", http.StatusUnauthorized)
	} else {
		http.Redirect(w, r, "/dashboard?username="+username, http.StatusFound)
	}
}

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

func route(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/signup_submit":
		signupSubmit(w, r)
	case "/login_submit":
		loginSubmit(w, r)
	case "/delete_user":
		database.DeleteUser(r.FormValue("username"))
	case "/timeout":
		http.Redirect(w, r, "/timeout", http.StatusGatewayTimeout)
	case "/":
		loadPage(w, r, "landing")
	default:
		title := strings.TrimLeft(r.URL.Path, "/")
		if _, exists := validEntryPoints[title]; exists {
			loadPage(w, r, title)
		} else {
			http.NotFound(w, r)
		}
	}
}

func main() {
	database.Init()
	http.HandleFunc("/", route)
	server := http.Server{
		Addr: ":8080",
	}
	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
