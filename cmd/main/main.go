// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"fmt"
	"github.com/CarlosACJ55/matcha/internal/database"
	"html/template"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	err := database.AuthenticateLogin(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		http.Redirect(w, r, "/login_fail.html", http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/dashboard.html", http.StatusOK)
}

func signupSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	password := r.FormValue("psw")
	if password != r.FormValue("psw-repeat") {
		http.Redirect(w, r, "signup_fail.html", http.StatusBadRequest)
		return
	}
	err := database.AddUser(r.FormValue("username"), r.FormValue("email"), password)
	if err != nil {
		http.Redirect(w, r, "login_fail.html", http.StatusUnauthorized)
	}
}

func loadPage(w http.ResponseWriter, fileName string) {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("ERROR when parsing file", err)
		return
	}
	err = t.ExecuteTemplate(w, fileName, nil)
	if err != nil {
		fmt.Println("ERROR when executing template", err)
	}
}

func handleFunction(w http.ResponseWriter, r *http.Request) {
	fmt.Printf((r.URL.Path) + "\n")
	switch r.URL.Path {
	case "/":
		loadPage(w, "landing.html")
	case "/login":
		loadPage(w, "login.html")
	case "/signup":
		loadPage(w, "signup.html")
	case "/login-submit":
		loginSubmit(w, r)
	case "/signup-submit":
		signupSubmit(w, r)
	case "/dashboard":
		loadPage(w, "dashboard.html")
	case "/settings":
		loadPage(w, "settings.html")
	case "/delete-user":
		database.DeleteUser(r.FormValue("username"))
	case "/timeout":
		if _, err := fmt.Fprint(w, "connection timed out"); err != nil {
			panic(err)
		}
	default:
		if _, err := fmt.Fprint(w, "nothing to see here"); err != nil {
			panic(err)
		}
	}
}

func main() {
	database.InitDB()
	http.HandleFunc("/", handleFunction)
	server := http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  10000000, // in ns
		WriteTimeout: 10000000, // in ns
	}
	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
