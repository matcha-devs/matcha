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

// TODO (@Alishah634) null pointer exception here, please fix before re-enabling.
//func executeSQLFile(filepath string) error {
//	file, err := os.Open(filepath)
//	if err != nil {
//		return err
//	}
//	defer func(file *os.File) {
//		err := file.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	}(file)
//	scanner := bufio.NewScanner(file)
//	var query strings.Builder
//	for scanner.Scan() {
//		line := scanner.Text()
//		if strings.HasPrefix(line, "--") { // Skip comments
//			continue
//		}
//		query.WriteString(line)
//		if strings.HasSuffix(line, ";") { // End of SQL statement
//			_, err := db.Exec(query.String())
//			if err != nil {
//				return err
//			}
//			query.Reset() // Reset query buffer for the next statement
//		}
//	}
//	fmt.Println("SQL file executed successfully")
//	if err := scanner.Err(); err != nil {
//		return err
//	}
//	return nil
//}

func InitDB() {
	password := os.Getenv("MYSQL_PASSWORD")
	rootDsn := "root:" + password + "@tcp(localhost:3306)/"

	// Connect to MySQL without specifying a database
	db, err := sql.Open("mysql", rootDsn)
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
	db, err = sql.Open("mysql", rootDsn+"matcha_db")
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
	s := string(text)
	_, err = db.Exec(s)
	if err != nil {
		log.Fatal("Error executing 'init.sql' - ", err)
	}
}

//func PrintUsersTable() {
//	rows, err := db.Query("SELECT * FROM users")
//	if err != nil {
//		fmt.Println("ERROR querying database", err)
//	}
//	defer func(rows *sql.Rows) {
//		err := rows.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	}(rows)
//
//	fmt.Println("==================================")
//	fmt.Println("id | username | email | password")
//	fmt.Println("---------------------------------")
//	for rows.Next() {
//		var user User
//		if err := rows.Scan(&user.id, &user.username, &user.email, &user.password); err != nil {
//			fmt.Printf("Error scanning row: %v\n", err)
//		}
//		fmt.Println(user.id, user.username, user.email, user.password)
//	}
//	fmt.Println("==================================")
//}

//func printOpenidTable() {
//	rows, err := db.Query("SELECT * FROM openid")
//	if err != nil {
//		fmt.Println("ERROR querying database", err)
//	}
//	defer func(rows *sql.Rows) {
//		err := rows.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	}(rows)
//
//	fmt.Println("========")
//	fmt.Println("Open id")
//	fmt.Println("-------")
//	for rows.Next() {
//		var id int˜˜
//		if err := rows.Scan(&id); err != nil {
//			fmt.Printf("Error scanning row: %v\n", err)
//		}
//		fmt.Println(id)
//	}
//	fmt.Println("========")
//}

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
	//checkOpenid
	var id int
	err := db.QueryRow("SELECT id FROM openid LIMIT 1").Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("There is no open ID")
	} else if err != nil {
		log.Fatalf("Error retrieving first row of openID table: %v", err)
	}

	query := "INSERT INTO users (username, email, password"
	if id != 0 { // if there is an open ID, assign it to the new user
		query += fmt.Sprintf(", id) VALUES (%s, %s, %s, %d)", username, email, password, id)
	} else {
		query += fmt.Sprintf(") VALUES (%s, %s, %s)", username, email, password)
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
