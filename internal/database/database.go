// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init() {
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
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS matcha_db")
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
	text, err := os.ReadFile("internal/query/init.sql")
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
		fmt.Println("There is no user. Running 'gen_users.sql' to create new users.")
		text, err := os.ReadFile("internal/query/gen_users.sql")
		if err != nil {
			log.Fatal("Error reading gen_users.sql file - ", err)
		}
		_, err = db.Exec(string(text))
		if err != nil {
			log.Fatal("Error executing 'gen_users.sql' - ", err)
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
	if dbPassword != password {
		err = errors.New("invalid password")
	}
	return err
}

func getOpenID() int {
	var id int
	err := db.QueryRow("SELECT id FROM openid LIMIT 1").Scan(&id)
	if id == 0 && !errors.Is(err, sql.ErrNoRows) {
		log.Fatalf("Error retrieving first row of openID table: %v", err)
	}
	return id
}

func AddUser(username string, email string, password string) error {
	var (
		query = "INSERT INTO users (username, email, password"
		id    = getOpenID()
	)
	if id == 0 { // if there is no open ID, assign a new id to the user.
		query += fmt.Sprintf(") VALUES (\"%s\", \"%s\", \"%s\")", username, email, password)
	} else { // Otherwise, reuse the open ID
		query += fmt.Sprintf(", id) VALUES (\"%s\", \"%s\", \"%s\", %d)", username, email, password, id)
	}
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error adding user - ", err)
	}
	fmt.Println("User Added Successfully")
	return err
}

func GetUserID(varName string, variable string) int {
	var id int
	err := db.QueryRow(fmt.Sprintf("SELECT id FROM users WHERE %s = ?", varName), variable).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error finding id using ", varName, " - ", err)
	}
	return id
}

func DeleteUser(id int) {
	_, err := db.Exec("INSERT INTO openid (id) VALUES(?)", id)
	if err != nil {
		log.Println("Error inserting openID ", id, " to the table - ", err)
	}
	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting the user id ", id, " - ", err)
	}
}

//TODO(): make a db.Close()
