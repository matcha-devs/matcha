// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	//use "go get -u github.com/go-sql-driver/mysql"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func login(w http.ResponseWriter) {
	var fileName = "login.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("ERROR when parsing file", err)
		return
	}

	err = t.ExecuteTemplate(w, fileName, nil)
	if err != nil {
		fmt.Println("ERROR when executing template", err)
		return
	}
}

func signup(w http.ResponseWriter) {
	var fileName = "signup.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("ERROR when parsing file", err)
		return
	}
	fmt.Println("hello")
	err = t.ExecuteTemplate(w, fileName, nil)
	if err != nil {
		fmt.Println("ERROR when executing template", err)
		return
	}
}

// searchUsername checks if the username exists in the database and returns the user ID
// func searchUsername(db *sql.DB, username string) (int, error) {
// 	var id int
// 	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return -1, nil // User not found
// 		}
// 		return -1, err // An error occurred
// 	}
// 	return id, nil
// }

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	db := InitDB() // Retrieve the singleton DB instance
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println("Username:", username)
	fmt.Println("Password:", password)

	if userValid, err := checkUser(db, username, password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Server error")
	} else if userValid {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Logged in successfully!")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Authentication failed")
	}
}

func checkUser(db *sql.DB, username, password string) (bool, error) {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	if err != nil {
		return false, err
	}
	return dbPassword == password, nil
}

func signupSubmit() {
}

func handleFunction(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		landing(w)
	case "/login":
		login(w)
	case "/login-submit":
		loginSubmit(w, r)
	case "/signup":
		signup(w)
	case "/signup-submit":
		signupSubmit(w, r)
	default:
		if _, err := fmt.Fprint(w, "nothing to see here"); err != nil {
			panic(err)
		}
	}
}

func landing(w http.ResponseWriter) {
	var fileName = "landing.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("ERROR when parsing file", err)
		return
	}

	err = t.ExecuteTemplate(w, fileName, nil)
	if err != nil {
		fmt.Println("ERROR when executing template", err)
		return
	}
}

func timeout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Timeout Attempt")
	time.Sleep(2 * time.Second)
	fmt.Fprint(w, "Did *not* timeout")
}

func main() {

	http.HandleFunc("/", handleFunction)
	http.HandleFunc("/timeout", timeout)
	http.HandleFunc("/login-submit", loginSubmit)

	server := http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  1000000,
		WriteTimeout: 1000000,
	}

	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
	//http protocol
	//http.ListenAndServeTLS("", "cert.pem", "key.pem", nil)
}
