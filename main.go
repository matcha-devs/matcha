// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
	if err := AuthenticateLogin(username, password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Authentication failed")
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Logged in successfully!")
}

func signupSubmit(w http.ResponseWriter, r *http.Request) {
	// db := InitDB() // Retrieve the singleton DB instance
	// Ask the user for their username, email, and password
	// Call the function addUser(db, username, email, password) this should add that instance of the user to the db
	// For debugging purposes, Print out the user's information and Print out the database's information, to confirm
	// that the user was added.
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("psw")
	repassword := r.FormValue("psw-repeat")
	//fmt.Println(email, username, password, repassword)
	fmt.Println("Signup Submit")
	if password == repassword {
		AddUser(username, email, password)
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
	default:
		if _, err := fmt.Fprint(w, "nothing to see here"); err != nil {
			panic(err)
		}
	}
}

func timeout(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Timeout Attempt")
	time.Sleep(2 * time.Second)
	fmt.Fprint(w, "Did *not* timeout")
	fmt.Println(w, "Did *not* timeout")
}

func main() {
	InitDB()
	printUsersTable()
	//Tester()

	http.HandleFunc("/", handleFunction)
	http.HandleFunc("/timeout", timeout)
	http.HandleFunc("/login-submit", loginSubmit)
	http.HandleFunc("/signup-submit", signupSubmit)
	server := http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  1000000, // in ns
		WriteTimeout: 1000000, // in ns
	}
	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
