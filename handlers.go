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

// Load to memory and generate all templates, panic if it fails.
var tmpl = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")))

// TODO(@FaaizMemonPurdue): This is an example of how go routines should be used, but we still need API call timeouts
// func routeWithTimeout(w http_server.ResponseWriter, r *http_server.Request) {
// 	ctx, cancel := context.WithTimeout(r.Context(), maxHandleTime)
// 	defer cancel()
// 	select {
// 	case <-ctx.Done():
// 		log.Println("Routing took longer than", maxHandleTime)
// 	default:
// 		start := time.Now()
// 		loadEntryPoint(w, r)
// 		log.Println("Routing done after", time.Since(start))
// 	}
// }

func loadPage(w http.ResponseWriter, r *http.Request, title string) {
	username := r.FormValue("username")
	id, err := matcha.database.GetUserID("username", username)
	if err != nil {
		log.Fatalf("Load Page Failed - can't get userID: ", err)
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
		log.Println("Login failed:", err)
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
		log.Println("Delete User failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		loadPage(w, r, "settings")
	} else {
		id, err := matcha.database.GetUserID("username", username)
		if err != nil {
			log.Fatalf("Delete User failed:", err)
		}
		err = matcha.database.DeleteUser(id)
		if err != nil {
			log.Println("Delete User failed:", err)
		}
	}
}
