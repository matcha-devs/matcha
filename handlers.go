package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/matcha-devs/matcha/internal"
)

var (
	// Load to memory and generate all resources, panic if it fails.
	tmpl         = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.gohtml")))
	publicServer = http.StripPrefix("/public", http.FileServer(http.Dir("public")))
	surfacePages = map[string]struct{}{
		"signup": {}, "signup-submit": {}, "signup-fail": {}, "login": {}, "login-submit": {}, "login-fail": {},
	}
)

func servePublicFile(w http.ResponseWriter, r *http.Request) {
	publicServer.ServeHTTP(w, r)
}

func loadIndex(w http.ResponseWriter, _ *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.go.html", nil)
	if err != nil {
		log.Println("Error executing template index.go.html -", err)
	}
}

func signupSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("signupSubmit")

	password := r.FormValue("psw")
	if password != r.FormValue("psw-repeat") {
		log.Println("Passwords didnt match.")
		http.Redirect(w, r, "/signup-fail", http.StatusSeeOther)
		return
	}

	// TODO(@seoyoungcho213): Validate user data, here or in the backend.
	username := r.FormValue("username")
	email := r.FormValue("email")

	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	err := matcha.database.AddUser(username, email, password)
	if err != nil {
		log.Println("Error adding user", username, "-", err)
		http.Redirect(w, r, "/signup-fail", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
	}
}

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("loginSubmit")

	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	id, err := matcha.database.AuthenticateLogin(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		log.Println("Login failed -", err)
		http.Redirect(w, r, "/login-fail", http.StatusSeeOther)
		return
	}
	http.SetCookie(
		w, &http.Cookie{
			Name:     "c_user_id",
			Value:    strconv.Itoa(id),
			Path:     "/",
			Expires:  time.Now().Add(20 * time.Minute),
			MaxAge:   20 * 60,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)
	http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	id, err := matcha.database.AuthenticateLogin(username, r.FormValue("password"))
	if err != nil {
		log.Println("User failed to validate delete request -", err)
		// TODO(@seoyoungcho213): Log the user out and immediately invalidate their cookie.
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	err = matcha.database.DeleteUser(id)
	if err != nil {
		log.Println("Delete User failed -", err)
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
	}

	// TODO(@seoyoungcho213) : Remove user id cookie
}

func checkLoginStatus(w http.ResponseWriter, r *http.Request) *internal.User {
	cookie, err := r.Cookie("c_user_id")
	if errors.Is(err, http.ErrNoCookie) {
		log.Println("Client has no session cookie -", err)
		http.Error(w, "Unauthorized login session.", http.StatusUnauthorized)
		return nil
	} else if err != nil {
		http.Error(w, "Unauthorized login session.", http.StatusUnauthorized)
		log.Println("Error getting session cookie -", err)
		return nil
	}
	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Println("Failed to convert user id -", err)
		http.Error(w, "Invalid login session.", http.StatusBadRequest)
		return nil
	} else if id < 1 {
		log.Println("Invalid user id -", err)
		http.Error(w, "Invalid login session.", http.StatusBadRequest)
		return nil
	}
	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	return matcha.database.GetUser(id)
}

func loadPage(w http.ResponseWriter, r *http.Request) {
	page := strings.TrimLeft(r.URL.Path, "/")
	log.Println("loading page {" + page + "}")
	var user *internal.User
	if _, exists := surfacePages[page]; !exists {
		user = checkLoginStatus(w, r)
		if user == nil {
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	err := tmpl.ExecuteTemplate(w, page+".gohtml", user)
	if err != nil {
		log.Println("Error executing template", page, "-", err)
	}
}
