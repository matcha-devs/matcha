package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/matcha-devs/matcha/structs"
)

var tmpl = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")))

func loadPage(w http.ResponseWriter, r *http.Request, title string) {
	username := r.FormValue("username")
	user := structs.User{
		ID:       deps.DB.GetUserID("username", username),
		Username: username,
		Email:    "test",
		Password: "test",
	}
	err := tmpl.ExecuteTemplate(w, title+".gohtml", user)
	if err != nil {
		log.Println("Error executing template-", err)
	}
}

func loadIndex(w http.ResponseWriter, r *http.Request) {
	loadPage(w, r, "index")
}

func loadEntryPoint(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	log.Println("Routing {" + path + "}")
	if _, exists := validEntryPoints[path]; !exists {
		log.Println("Not a valid entry point:", path)
		http.NotFound(w, r)
	}
	loadPage(w, r, path)
}

func signupSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	password := r.FormValue("psw")
	if password != r.FormValue("psw-repeat") {
		log.Println("Passwords didnt match.")
		http.Redirect(w, r, "/signup-fail", http.StatusSeeOther)
		return
	}
	username := r.FormValue("username")
	err := deps.DB.AddUser(username, r.FormValue("email"), password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "signup-fail")
	} else {
		loadPage(w, r, "dashboard")
	}
}

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	username := r.FormValue("username")
	err := deps.DB.AuthenticateLogin(username, r.FormValue("password"))
	if err != nil {
		log.Println("Login failed:", err)
		http.Redirect(w, r, "/login-fail", http.StatusSeeOther)
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
	err := deps.DB.AuthenticateLogin(username, password)
	if err != nil {
		log.Println("Delete User failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "settings")
	} else {
		id := deps.DB.GetUserID("username", username)
		err := deps.DB.DeleteUser(id)
		if err != nil {
			log.Println("Delete User failed:", err)
		}
	}
}
