package main

import (
	"log"
	"net/http"

	"github.com/matcha-devs/matcha/internal/sql"
)

func loadPage(w http.ResponseWriter, r *http.Request, title string) {
	username := r.FormValue("username")
	user := sql.User{
		ID:       deps.db.GetUserID("username", username),
		Username: username,
		Email:    "test",
		Password: "test",
	}
	err := tmpl.ExecuteTemplate(w, title+".gohtml", user)
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
	err := deps.db.AddUser(username, r.FormValue("email"), password)
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
	err := deps.db.AuthenticateLogin(username, r.FormValue("password"))
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
	err := deps.db.AuthenticateLogin(username, password)
	if err != nil {
		log.Println("Delete User failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "settings")
	} else {
		id := deps.db.GetUserID("username", username)
		err := deps.db.DeleteUser(id)
		if err != nil {
			log.Println("Delete User failed:", err)
		}
	}
}
