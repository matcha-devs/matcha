// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	//use "go get -u github.com/go-sql-driver/mysql"
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

func searchUsername(username string) int {
	userDB := GetUsers()
	for _, user := range userDB {
		if user.Name == username {
			return user.ID
		}
	}
	return -1
}

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if id := searchUsername(username); id == -1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "User not found!")
	} else if GetPassword(id) == password {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Logged in successfully!")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Wrong Password")
	}
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
	//fmt.Println(os.Getenv("MYSQL_PASSWORD"))
	//pswd := os.Getenv("MYSQL_PASSWORD")
	//db, err := sql.Open("mysql", "root:"+pswd+"@tcp(localhost:3306)/userdb")
	//
	//if err != nil {
	//	fmt.Println("ERROR when opening database connection", err)
	//	panic(err.Error())
	//}
	//defer db.Close()
	//fmt.Println("Successfully opened database connection")

	http.HandleFunc("/", handleFunction)
	http.HandleFunc("/timeout", timeout)

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
