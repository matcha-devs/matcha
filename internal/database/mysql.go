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
	"github.com/matcha-devs/matcha/internal"
	"golang.org/x/crypto/bcrypt"
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
		log.Fatalln("Error reading init.sql file -", err)
	}

	// Connect to MySQL root.
	db, err := sql.Open("mysql", mysql.rootDsn)
	if err != nil {
		log.Fatalln("Error opening MySQL rootDsn -", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		log.Fatalln("Error connecting to MySQL rootDsn -", err)
	}

	// Create the database if it does not exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		log.Fatalln("Error creating database -", err)
	}
	if err := db.Close(); err != nil {
		log.Fatalln("Error closing database -", err)
	}

	// Connect to the database
	db, err = sql.Open("mysql", mysql.rootDsn+dbName+"?multiStatements=true")
	if err != nil {
		log.Fatalln("Error opening database -", err)
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		log.Fatalln("Error connecting to database -", err)
	}

	// Run 'init.sql' script.
	_, err = db.Exec(string(initScript))
	if err != nil {
		log.Fatalln("Error executing 'init.sql' -", err)
	}
	if err := db.Close(); err != nil {
		log.Fatalln("Error closing database -", err)
	}
	return &mysql
}

func (db *MySQLDatabase) Open() (err error) {
	db.underlyingDB, err = sql.Open("mysql", db.rootDsn+db.dbName+"?parseTime=true")
	if err != nil {
		log.Fatalln("Error opening database -", err)
		return err
	}
	log.Println("MySQL Database connecting to", db.rootDsn[strings.Index(db.rootDsn, "@"):]+db.dbName, "🫡")
	if err := db.underlyingDB.Ping(); err != nil {
		_ = db.underlyingDB.Close()
		log.Fatalln("Error connecting to database -", err)
		return err
	}
	return err
}

func (db *MySQLDatabase) Close() (err error) {
	if err := db.underlyingDB.Close(); err != nil {
		log.Println("underlying database close failure -", err)
		return err
	}
	log.Println("MySQL database has closed 👋🏽")
	return nil
}

func (db *MySQLDatabase) AuthenticateLogin(username string, password string) (id int, err error) {
	var expectedPassword string
	err = db.underlyingDB.QueryRow(
		"SELECT id, password FROM users WHERE BINARY username = ?", username,
	).Scan(&id, &expectedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("invalid username")
	}
	err = bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(password))
	if err != nil {
		return 0, errors.New("invalid password")
	}
	return id, err
}

func (db *MySQLDatabase) GetUser(id int) *internal.User {
	var user = &internal.User{}
	err := db.underlyingDB.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedOn,
	)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("No user with ID:", id, "-", err)
		return nil
	} else if err != nil {
		log.Println("Failed to query users for ID:", id, "-", err)
		return nil
	} else if !(user.ID.Valid &&
		user.Username.Valid && user.Email.Valid && user.Password.Valid && user.CreatedOn.Valid) {
		log.Println("Malformed user with ID:", id, "-", user)
		return nil
	}
	return user
}

func (db *MySQLDatabase) getOpenID() int {
	var id int
	err := db.underlyingDB.QueryRow("SELECT id FROM openid LIMIT 1").Scan(&id)
	if id == 0 && !errors.Is(err, sql.ErrNoRows) {
		log.Fatalln("Error retrieving first row of openID table - ", err)
	}
	return id
}

func (db *MySQLDatabase) AddUser(username string, email string, password string) (err error) {
	// TODO(@seoyoungcho213): For efficiency, we might be able to return the new id here with only a single query?
	var (
		query = "INSERT INTO users (username, email, password"
		id    = db.getOpenID()
	)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password -", err)
		return err
	}
	if id == 0 { // if there is no open ID, assign a new id to the user.
		query += fmt.Sprintf(") VALUES (\"%s\", \"%s\", \"%s\")", username, email, hashedPassword)
	} else { // Otherwise, reuse the open ID
		query += fmt.Sprintf(", id) VALUES (\"%s\", \"%s\", \"%s\", %d)", username, email, hashedPassword, id)
	}
	_, err = db.underlyingDB.Exec(query)
	if err != nil {
		log.Println("Error adding user -", err)
	}
	return err
}

func (db *MySQLDatabase) GetUserID(varName string, variable string) int {
	var id int
	err := db.underlyingDB.QueryRow(
		fmt.Sprintf("SELECT id FROM users WHERE BINARY %s = ?", varName), variable,
	).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error querying users for", varName+":"+variable, "-", err)
		return 0
	}
	return id
}

func (db *MySQLDatabase) DeleteUser(id int) (err error) {
	if _, err = db.underlyingDB.Exec("INSERT INTO openid (id) VALUES(?)", id); err != nil {
		log.Println("Error inserting openID", id, " to the table -", err)
		return err
	}
	if _, err = db.underlyingDB.Exec("DELETE FROM users WHERE BINARY id = ?", id); err != nil {
		log.Println("Error deleting the user id", id, " -", err)
		return err
	}
	return nil
}
