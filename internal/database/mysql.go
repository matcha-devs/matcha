// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/matcha-devs/matcha/internal"
	"golang.org/x/crypto/bcrypt"
)

type MySQLDatabase struct {
	rootDSN      string
	dbName       string
	underlyingDB *sql.DB
}

func New(dbName string, username string, password string) (mysql *MySQLDatabase) {
	// Initialized struct to be returned.
	mysql = &MySQLDatabase{
		rootDSN:      username + ":" + password + "@tcp(localhost:3306)/",
		dbName:       dbName,
		underlyingDB: nil,
	}

	// Open a separate connection to the root DSN and create the database if it does not exist
	initDB, err := sql.Open("mysql", mysql.rootDSN+"?multiStatements=true")
	if err != nil {
		log.Fatalln("Error opening MySQL root DSN -", err)
	}
	if err = initDB.Ping(); err != nil {
		if err = initDB.Close(); err != nil {
			log.Fatalln("Error closing broken MySQL root DSN -", err)
		}
		log.Fatalln("Error connecting to MySQL root DSN -", err)
	}
	_, err = initDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatalln("Error creating database -", err)
	}

	// Move to the database and run 'init_tables.sql' script.
	_, err = initDB.Exec("USE " + dbName)
	if err != nil {
		log.Fatalln("Error using", dbName, "-", err)
	}
	initScript, err := os.ReadFile("internal/database/queries/init_tables.sql")
	if err != nil {
		log.Fatalln("Error reading init_tables.sql file -", err)
	}
	if _, err = initDB.Exec(string(initScript)); err != nil {
		log.Fatalln("Error executing init_tables.sql -", err)
	}
	if err := initDB.Close(); err != nil {
		log.Fatalln("Error closing init DB -", err)
	}
	return
}

func (db *MySQLDatabase) Open() (err error) {
	if db.underlyingDB, err = sql.Open("mysql", db.rootDSN+db.dbName+"?parseTime=true"); err != nil {
		log.Println("Error opening database -", err)
		return
	}
	log.Println("MySQL Database connecting to", db.rootDSN[strings.Index(db.rootDSN, "@"):]+db.dbName, "🫡")
	if err = db.underlyingDB.Ping(); err != nil {
		if err = db.underlyingDB.Close(); err != nil {
			log.Println("Error closing broken database -", err)
		}
		log.Println("Error connecting to database -", err)
	}
	return
}

func (db *MySQLDatabase) Close() (err error) {
	if err = db.underlyingDB.Close(); err != nil {
		log.Println("underlying database close failure -", err)
	} else {
		log.Println("MySQL database has closed 👋🏽")
	}
	return
}

func (db *MySQLDatabase) AuthenticateLogin(email, password string) (id uint64, err error) {
	var hash []byte
	err = db.underlyingDB.QueryRow("SELECT id, password FROM users WHERE BINARY email = ?", email).Scan(
		&id, &hash,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("invalid email")
	}
	if err = bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return 0, errors.New("invalid password")
	}
	return
}

func (db *MySQLDatabase) GetUser(id uint64) (user *internal.User) {
	user = &internal.User{}
	err := db.underlyingDB.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(
		&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.Password, &user.DateOfBirth,
		&user.CreatedOn)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("No user with ID:", id, "-", err)
		return nil
	} else if err != nil {
		log.Println("Failed to query users for ID:", id, "-", err)
		return nil
	} else if !user.IsValid() {
		log.Println("Malformed user with ID:", id, "-", user)
		return nil
	}
	return
}

func (db *MySQLDatabase) getOpenID() (id uint64, err error) {
	err = db.underlyingDB.QueryRow("SELECT id FROM openid LIMIT 1").Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Failed to query for an openid -", err)
		return 0, errors.New("invalid openid")
	}
	return
}

func (db *MySQLDatabase) AddUser(firstName, middleName, lastName, email, password, dateOfBirth string) (
	id uint64, err error) {
	if len(firstName) == 0 || len(lastName) == 0 || len(email) == 0 || len(password) == 0 || len(dateOfBirth) == 0 {
		return 0, errors.New("empty fields")
	}
	query := "INSERT INTO users (first_name, middle_name, last_name, email, password, date_of_birth"
	if openID, err := db.getOpenID(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error getting open id -", err)
		return id, errors.New("internal server error")
	} else if openID == 0 {
		query += `) VALUES (?, ?, ?, ?, ?, ?)`
	} else {
		log.Println("Re-using open id:", openID, "for {"+email+"}")
		id = openID
		query += `, id) VALUES (?, ?, ?, ?, ?, ?, ` + strconv.FormatUint(id, 10) + `)`
		if _, err = db.underlyingDB.Exec("DELETE FROM openid WHERE id = ?", id); err != nil {
			log.Println("Error deleting open id -", err)
			return id, errors.New("internal server error")
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password -", err)
		return id, errors.New("internal server error")
	}
	result, err := db.underlyingDB.Exec(query, firstName, middleName, lastName, email, hashedPassword, dateOfBirth)
	if err != nil {
		log.Println("Error adding user -", err)
		return 0, errors.New("internal server error")
	}
	if id == 0 {
		userid, err := result.LastInsertId()
		if err != nil {
			log.Println("Error getting user ID -", err)
			return 0, errors.New("internal server error")
		}
		id = uint64(userid)
		log.Println("All existing IDs in use, assigning new ID:", id, "to {"+email+"}")
	}
	return
}

func (db *MySQLDatabase) GetUserID(email string) (id uint64) {
	// TODO(@seoyoungcho213): might not use this anymore cuz of cookie
	if err := db.underlyingDB.QueryRow(
		"SELECT id FROM users WHERE BINARY email = ?", email,
	).Scan(&id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error querying users for email:"+email, "-", err)
		return 0
	}
	return id
}

func (db *MySQLDatabase) DeleteUser(id uint64) (err error) {
	if _, err = db.underlyingDB.Exec("INSERT INTO openid (id) VALUES(?)", id); err != nil {
		log.Println("Error inserting openID", id, " to the table -", err)
		return errors.New("internal server error")
	}
	if _, err = db.underlyingDB.Exec("DELETE FROM users WHERE BINARY id = ?", id); err != nil {
		log.Println("Error deleting the user id", id, " -", err)
		return errors.New("internal server error")
	}
	return err
}
