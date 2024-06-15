// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package mySQL

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB 
}

func Open(database_name string, queries_path string) *Database {
	password := os.Getenv("MYSQL_PASSWORD")
	rootDsn := "root:" + password + "@tcp(localhost:3306)/"
	// Connect to MySQL without specifying matchaDB
	db, err := sql.Open("mysql", rootDsn)
	if err != nil {
		log.Fatal("Error opening matchaDB-", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to MySQL-", err)
	}

	// Create the matcha_db if it does not exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + database_name)
	if err != nil {
		log.Fatal("Error opening Database-", err)
	}
	if err := db.Close(); err != nil {
		log.Fatal("Error closing Database-", err)
	}

	// Connect to matcha_db to run 'init.sql' script
	db, err = sql.Open("mysql", rootDsn+database_name+"?multiStatements=true")
	if err != nil {
		log.Fatal("Error opening Database-", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to Database-", err)
	}
	text, err := os.ReadFile(queries_path+"init.sql")
	if err != nil {
		log.Fatal("Error reading init.sql file-", err)
	}
	_, err = db.Exec(string(text))
	if err != nil {
		log.Fatal("Error executing 'init.sql'-", err)
	}

	// If there is no user, then make test users.
	var userCount int
	err = db.QueryRow("SELECT COUNT(*) AS user_count FROM users").Scan(&userCount)
	if err != nil {
		log.Fatal(err)
	}
	if userCount == 0 {
		fmt.Println("There is no user. Running 'gen_users.sql' to create new users.")
		text, err := os.ReadFile(queries_path+"gen_users.sql")
		if err != nil {
			log.Fatal("Error reading gen_users.sql file-", err)
		}
		_, err = db.Exec(string(text))
		if err != nil {
			log.Fatal("Error executing 'gen_users.sql'-", err)
		}
	}

	// Re-open matchaDB for the security purpose
	if err := db.Close(); err != nil {
		log.Println("Error closing Database-", err)
	}
	db, err = sql.Open("mysql", rootDsn+database_name)
	if err != nil {
		log.Fatal("Error opening Database-", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to Database-", err)
	}
	fmt.Println("HERE!!")
	return &Database{db}
}

func (matcha Database) AuthenticateLogin(username string, password string) error {
	var dbPassword string
	err := matcha.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	if dbPassword != password {
		err = errors.New("invalid password")
	}
	return err
}

func (matcha Database) getOpenID() int {
	var id int
	err := matcha.db.QueryRow("SELECT id FROM openid LIMIT 1").Scan(&id)
	if id == 0 && !errors.Is(err, sql.ErrNoRows) {
		log.Fatalf("Error retrieving first row of openID table: %v", err)
	}
	return id
}

func (matcha Database) AddUser(username string, email string, password string) error {
	var (
		query = "INSERT INTO users (username, email, password"
		id    = matcha.getOpenID()
	)
	if id == 0 { // if there is no open ID, assign a new id to the user.
		query += fmt.Sprintf(") VALUES (\"%s\", \"%s\", \"%s\")", username, email, password)
	} else { // Otherwise, reuse the open ID
		query += fmt.Sprintf(", id) VALUES (\"%s\", \"%s\", \"%s\", %d)", username, email, password, id)
	}
	_, err := matcha.db.Exec(query)
	if err != nil {
		log.Fatal("Error adding user-", err)
	}
	fmt.Println("User Added Successfully")
	return err
}

func (matcha Database) GetUserID(varName string, variable string) int {
	var id int
	err := matcha.db.QueryRow(fmt.Sprintf("SELECT id FROM users WHERE %s = ?", varName), variable).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error finding id using ", varName, "-", err)
	}
	return id
}

func (matcha Database) DeleteUser(id int) error {
	_, err := matcha.db.Exec("INSERT INTO openid (id) VALUES(?)", id)
	if err != nil {
		log.Println("Error inserting openID ", id, " to the table - ", err)
		return err
	}
	_, err = matcha.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting the user id ", id, " - ", err)
		return err
	}
	return nil
}

func (matcha Database) Close() error {
	if err := matcha.db.Close(); err != nil {
		log.Println("mysqldb close failure:", err)
		return err
	}
	return nil
}
