package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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
	var username, email, password string
	id := getCookieHandler(w, r, "c_user")
	if id == 0 {
		username = "user_unknown"
		email = "user_unknown"
		password = "user_unknown"
	} else {
		username, email, password, err = matcha.database.GetUserInfo(id)
		if err != nil {
			http.Error(w, "Error occurred from getting user info", http.StatusInternalServerError)
			return
		}
	}
	user := structs.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
	err = tmpl.ExecuteTemplate(w, title+".gohtml", user)
	if err != nil {
		log.Println("Error executing template", title, "-", err)
	}
}

func loadEntryPoint(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
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
		return
	}
	cookie, err := r.Cookie("c_user")
	if err != nil {
		fmt.Println("cookie was not found")
		id, err := matcha.database.GetUserID("username", username)
		if err != nil {
			fmt.Println("getting user id failed -", err)
		}
		cookie = &http.Cookie{
			Name:     "c_user",
			Value:    strconv.Itoa(id),
			Expires:  time.Now().Add(20 * time.Minute),
			MaxAge: 20 * time.Minute
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
	}
	loadPage(w, r, "dashboard")

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

func getCookieHandler(w http.ResponseWriter, r *http.Request, cookie_name string) int {
	value, err := cookies.Read(r, cookie_name)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		case errors.Is(err, cookies.ErrInvalidValue):
			http.Error(w, "invalid cookie", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return 0
	}
	return strconv.Atoi(value)
}