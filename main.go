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
	fmt.Println("hello")
	err = t.ExecuteTemplate(w, fileName, nil)
	if err != nil {
		fmt.Println("ERROR when executing template", err)
		return
	}
}

// searchUsername checks if the username exists in the database and returns the user ID
func searchUsername(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil // User not found
		}
		return -1, err // An error occurred
	}
	return id, nil
}

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	db := GetDB() // Retrieve the singleton DB instance
	username := r.FormValue("username")
	password := r.FormValue("password")

	id, err := searchUsername(db, username)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if id == -1 {
		http.Error(w, "User not found!", http.StatusNotFound)
		return
	}

	storedPassword, err := GetPassword(db, id)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if storedPassword == password {
		fmt.Fprint(w, "Logged in successfully!")
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
	}
}

func handleFunction(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		if _, err := fmt.Fprint(w, "<h1>Welcome to Menthol!</h1>"); err != nil {
			panic(err)
		}
	case "/login":
		login(w)
	case "/login-submit":
		loginSubmit(w, r)
	default:
		if _, err := fmt.Fprint(w, "nothing to see here"); err != nil {
			panic(err)
		}
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

	server := http.Server{
		Addr:         "",
		Handler:      nil,
		ReadTimeout:  1000,
		WriteTimeout: 1000,
	}

	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
	//http protocol
	//http.ListenAndServeTLS("", "cert.pem", "key.pem", nil)
}
