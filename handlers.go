package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/matcha-devs/matcha/internal/structs"
)

var (
	// Load to memory and generate all resources, panic if it fails.
	tmpl         = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")))
	publicServer = http.StripPrefix("/public", http.FileServer(http.Dir("public")))
)

// TODO(@FaaizMemonPurdue): Add API call timeouts.

func loadPage(w http.ResponseWriter, r *http.Request, title string) {
	username := r.FormValue("username")
	id, err := matcha.database.GetUserID("username", username)
	if err != nil {
		log.Println("Load Page Failed - can't get userID: ", err)
	}
	user := structs.User{
		ID:        id,
		Username:  username,
		Email:     "test",
		Password:  "test",
		CreatedAt: time.Now(),
	}
	err = tmpl.ExecuteTemplate(w, title+".gohtml", user)
	if err != nil {
		log.Println("Error executing template -", err)
	}
}

func loadEntryPoint(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	log.Println("Routing {" + path + "}")
	if _, exists := validEntryPoints[path]; !exists {
		log.Println("Not a valid entry point -", path)
		http.NotFound(w, r)
	}
	loadPage(w, r, path)
}

func loadIndex(w http.ResponseWriter, r *http.Request) {
	loadPage(w, r, "index")
}

func servePublicFile(w http.ResponseWriter, r *http.Request) {
	publicServer.ServeHTTP(w, r)
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
	err := matcha.database.AddUser(username, r.FormValue("email"), password)
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
	err := matcha.database.AuthenticateLogin(username, r.FormValue("password"))
	if err != nil {
		log.Println("Login failed -", err)
		http.Redirect(w, r, "/login-fail", http.StatusSeeOther)
	} else {
		loadPage(w, r, "dashboard")
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "invalid request", http.StatusSeeOther)
		loadPage(w, r, "/")
		return
	}
	username := r.FormValue("username")
	err := matcha.database.AuthenticateLogin(username, r.FormValue("password"))
	if err != nil {
		log.Println("Delete User failed -", err)
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "settings")
	} else {
		id, err := matcha.database.GetUserID("username", username)
		if err != nil {
			log.Println("Delete User failed -", err)
		}
		err = matcha.database.DeleteUser(id)
		if err != nil {
			log.Println("Delete User failed -", err)
		}
	}
}
