// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	password := os.Getenv("MYSQL_PASSWORD")
	rootDsn := "root:" + password + "@tcp(localhost:3306)/"
	var err error
	// Connect to MySQL without specifying a database
	db, err = sql.Open("mysql", rootDsn)
	if err != nil {
		log.Fatal("Error opening database - ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to MySQL - ", err)
	}

	// Create the matcha_db if it does not exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS matcha_db;")
	if err != nil {
		log.Fatal("Error opening matcha_db - ", err)
	}
	if err := db.Close(); err != nil {
		return
	}

	// Connect to matcha_db to run 'init.sql' script
	db, err = sql.Open("mysql", rootDsn+"matcha_db?multiStatements=true")
	if err != nil {
		log.Fatal("Error opening matcha_db - ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to matcha_db - ", err)
	}
	text, err := os.ReadFile("internal/database/init.sql")
	if err != nil {
		log.Fatal("Error reading init.sql file - ", err)
	}
	_, err = db.Exec(string(text))
	if err != nil {
		log.Fatal("Error executing 'init.sql' - ", err)
	}

	// If there is no user, then make test users.
	var userCount int
	err = db.QueryRow("SELECT COUNT(*) AS user_count FROM users").Scan(&userCount)
	if err != nil {
		log.Fatal(err)
	}
	if userCount == 0 {
		fmt.Println("There is no user. Running 'testusers.sql' to create new users.")
		text, err := os.ReadFile("internal/database/testusers.sql")
		if err != nil {
			log.Fatal("Error reading testusers.sql file - ", err)
		}
		_, err = db.Exec(string(text))
		if err != nil {
			log.Fatal("Error executing 'testusers.sql' - ", err)
		}
	}

	// Re-open the database for the security purpose
	if err := db.Close(); err != nil {
		return
	}
	db, err = sql.Open("mysql", rootDsn+"matcha_db")
	if err != nil {
		log.Fatal("Error opening matcha_db - ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to matcha_db - ", err)
	}
}

func AuthenticateLogin(username, password string) error {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	fmt.Println("DB Password:", dbPassword)
	fmt.Println("err:", err)

	if err != nil {
		return err
	} else if dbPassword != password {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func AddUser(username string, email string, password string) error {
	// check number of Openid
	var id int
	err := db.QueryRow("SELECT COUNT(*) AS id_count FROM openid").Scan(&id)
	if err != nil {
		log.Fatalf("Error querying openID: %v", err)
	}
	if id != 0 { // if there is open ID, assign it to id.
		err = db.QueryRow("SELECT id FROM openid LIMIT 1").Scan(&id)
		if err != nil {
			log.Fatalf("Error retrieving first row of openID table: %v", err)
		}
	}

	query := "INSERT INTO users (username, email, password"
	if id == 0 { // if there is no open ID, don't assign id to the new user.
		query += fmt.Sprintf(") VALUES (\"%s\", \"%s\", \"%s\");", username, email, password)
	} else {
		query += fmt.Sprintf(", id) VALUES (\"%s\", \"%s\", \"%s\", %d);", username, email, password, id)
	}

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Error adding user : %v\n", err)
	}

	fmt.Println("User Added Successfully")
	return err
}

func FindID(varName string, variable string) int {
	var id int
	if varName == "username" {
		row := db.QueryRow("SELECT id FROM users WHERE username = ?;", variable)
		err := row.Scan(&id)
		if err != nil {
			log.Fatalf("Error finding id using username: %v", err)
		}
	} else if varName == "email" {
		row := db.QueryRow("SELECT id FROM users WHERE email = ?;", variable)
		err := row.Scan(&id)
		if err != nil {
			log.Fatalf("Error finding id using email: %v", err)
		}
	}
	return id
}

func DeleteUser(username string) {
	// add id to openID table
	id := FindID("username", username)
	_, err := db.Exec("INSERT INTO openid (id) VALUES(?);", id)
	if err != nil {
		log.Fatalf("Error inserting openID to the table: %v", err)
	}

	_, err = db.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		log.Fatalf("Error deleting the user: %v", err)
	}
}
