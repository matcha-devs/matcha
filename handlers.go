package main

import (
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/matcha-devs/matcha/internal"
)

var (
	// Load to memory and generate all resources, panic if it fails.
	tmpl         = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.go.html")))
	publicServer = http.StripPrefix("/public", http.FileServer(http.Dir("public")))
	surfacePages = map[string]struct{}{"signup": {}, "login": {}}
)

func getPublic(w http.ResponseWriter, r *http.Request) {
	publicServer.ServeHTTP(w, r)
}

func getIndex(w http.ResponseWriter, _ *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.go.html", nil)
	if err != nil {
		log.Println("Error executing template index.go.html -", err)
	}
}

func setSessionCookie(w http.ResponseWriter, id int) {
	log.Println("Issued cookie for id:", id)
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
}

func postSignup(w http.ResponseWriter, r *http.Request) {
	firstname := r.FormValue("firstname")
	// TODO(@seoyoungcho213): Validate user data way better here.
	middlename := r.FormValue("middlename")
	lastname := r.FormValue("lastname")
	email := r.FormValue("email")
	// TODO(@seoyoungcho213): read about mail.ParseAddress' return type and use it to add first+last name formatting too
	if _, err := mail.ParseAddress("<" + email + ">"); err != nil {
		log.Println("Error adding user {"+email+"} to database -", err)
		if _, err := io.WriteString(w, err.Error()); err != nil {
			log.Println("Error writing signup email format error -", err)
		}
		return
	}
	password := r.FormValue("psw")
	if password != r.FormValue("psw-repeat") {
		log.Println("Passwords didnt match.")
		if _, err := io.WriteString(w, "passwords do not match"); err != nil {
			log.Println("Error writing signup failure -", err)
		}
		return
	}
	birthdate := r.FormValue("birthdate")
	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	id, err := matcha.database.AddUser(firstname, middlename, lastname, email, password, birthdate)
	if err != nil {
		log.Println("Error adding user {"+email+"} to database -", err)
		if _, err := io.WriteString(w, err.Error()); err != nil {
			log.Println("Error writing server error -", err)
		}
		return
	}
	setSessionCookie(w, id)
	w.Header().Set("HX-Redirect", "/dashboard")
}

func postLogout(w http.ResponseWriter, _ *http.Request) {
	setSessionCookie(w, 0)
	w.Header().Set("HX-Redirect", "/")
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	id, err := matcha.database.AuthenticateLogin(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		log.Println("Login failed -", err)
		if _, err := io.WriteString(w, err.Error()); err != nil {
			log.Println("Error writing login failure -", err)
		}
		return
	}
	setSessionCookie(w, id)
	w.Header().Set("HX-Redirect", "/dashboard")
}

func postDeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	username := r.FormValue("username")
	id, err := matcha.database.AuthenticateLogin(username, r.FormValue("password"))
	if err != nil {
		log.Println("User failed to validate delete request -", err)
		postLogout(w, r)
		return
	}

	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	err = matcha.database.DeleteUser(id)
	if err != nil {
		log.Println("Delete User failed -", err)
		if _, err = io.WriteString(w, "internal server error"); err != nil {
			log.Println("Error writing deleted user internal server error -", err)
			return
		}
	}
	postLogout(w, r)
}

func checkLoginStatus(w http.ResponseWriter, r *http.Request) (user *internal.User) {
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
		log.Println("Invalid user id:", id)
		http.Error(w, "Invalid login session.", http.StatusBadRequest)
		return nil
	}
	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	return matcha.database.GetUser(id)
}

func getPage(w http.ResponseWriter, r *http.Request) {
	page := strings.TrimLeft(r.URL.Path, "/")
	var user *internal.User
	if _, exists := surfacePages[page]; !exists {
		user = checkLoginStatus(w, r)
		if user == nil {
			return
		}
	}
	if err := tmpl.ExecuteTemplate(w, page+".go.html", user); err != nil {
		log.Println("Error executing template", page, "-", err)
	}
}
