// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDatabase struct {
	rootDsn      string
	dbName       string
	underlyingDB *sql.DB
}

func New(dbName string, username string, password string) *MySQLDatabase {
	mysql := MySQLDatabase{
		rootDsn:      username + ":" + password + "@tcp(localhost:3306)/",
		dbName:       dbName,
		underlyingDB: nil,
	}
	initScript, err := os.ReadFile("internal/database/queries/init.sql")
	if err != nil {
		log.Fatal("Error reading init.sql file -", err)
	}

	// Connect to MySQL root.
	db, err := sql.Open("mysql", mysql.rootDsn)
	if err != nil {
		log.Fatal("Error opening MySQL rootDsn -", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to MySQL rootDsn -", err)
	}

	// Create the database if it does not exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatal("Error creating database-", err)
	}
	if err := db.Close(); err != nil {
		log.Fatal("Error closing database -", err)
	}

	// Connect to the database
	db, err = sql.Open("mysql", mysql.rootDsn+dbName+"?multiStatements=true")
	if err != nil {
		log.Fatal("Error opening database -", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to database -", err)
	}

	// Run 'init.sql' script.
	_, err = db.Exec(string(initScript))
	if err != nil {
		log.Fatal("Error executing 'init.sql'-", err)
	}
	if err := db.Close(); err != nil {
		log.Println("Error closing database -", err)
	}
	return &mysql
}

func (db *MySQLDatabase) Open() error {
	var err error
	db.underlyingDB, err = sql.Open("mysql", db.rootDsn+db.dbName)
	if err != nil {
		log.Println("Error opening database -", err)
		return err
	}
	log.Println("MySQL Database connecting to", db.rootDsn[strings.Index(db.rootDsn, "@"):]+db.dbName, "🫡")
	if err := db.underlyingDB.Ping(); err != nil {
		log.Println("Error connecting to database -", err)
		return err
	}
	return err
}

func (db *MySQLDatabase) Close() error {
	if err := db.underlyingDB.Close(); err != nil {
		log.Println("underlying database close failure -", err)
		return err
	}
	log.Println("MySQL database has closed 👋🏽")
	return nil
}

func (db *MySQLDatabase) AuthenticateLogin(username string, password string) error {
	var dbPassword string
	err := db.underlyingDB.QueryRow(
		"SELECT password FROM users WHERE BINARY username = ?", username,
	).Scan(&dbPassword)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("invalid username")
	} else if dbPassword != password {
		err = errors.New("invalid password")
	}
	return err
}

func (db *MySQLDatabase) getOpenID() int {
	var id int
	err := db.underlyingDB.QueryRow("SELECT id FROM openid LIMIT 1").Scan(&id)
	if id == 0 && !errors.Is(err, sql.ErrNoRows) {
		log.Fatalf("Error retrieving first row of openID table: %v", err)
	}
	return id
}

func (db *MySQLDatabase) AddUser(username string, email string, password string) error {
	var (
		query = "INSERT INTO users (username, email, password"
		id    = db.getOpenID()
	)
	if id == 0 { // if there is no open ID, assign a new id to the user.
		query += fmt.Sprintf(") VALUES (\"%s\", \"%s\", \"%s\")", username, email, password)
	} else { // Otherwise, reuse the open ID
		query += fmt.Sprintf(", id) VALUES (\"%s\", \"%s\", \"%s\", %d)", username, email, password, id)
	}
	_, err := db.underlyingDB.Exec(query)
	if err != nil {
		log.Fatal("Error adding user -", err)
	}
	fmt.Println("User Added Successfully")
	return err
}

func (db *MySQLDatabase) GetUserID(varName string, variable string) (int, error) {
	var id int
	err := db.underlyingDB.QueryRow(
		fmt.Sprintf("SELECT id FROM users WHERE BINARY %s = ?", varName), variable,
	).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error finding id using", varName, "-", err)
	}
	return id, err
}

func (db *MySQLDatabase) DeleteUser(id int) error {
	_, err := db.underlyingDB.Exec("INSERT INTO openid (id) VALUES(?)", id)
	if err != nil {
		log.Println("Error inserting openID", id, " to the table -", err)
		return err
	}
	_, err = db.underlyingDB.Exec("DELETE FROM users WHERE BINARY id = ?", id)
	if err != nil {
		log.Println("Error deleting the user id", id, " -", err)
		return err
	}
	return nil
}
