package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"sync"

	// The following imports to use the First database instance:
	"bufio"
	"fmt"
	"strings"
)

var once sync.Once
var db *sql.DB

// InitDB returns a singleton database instance
func InitDB() {
	once.Do(func() {
		var err error
		pswd := os.Getenv("MYSQL_PASSWORD") // Ensure this environment variable is set
		dsn := "root:" + pswd + "@tcp(127.0.0.1:3306)/userdb"
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}
		if err = db.Ping(); err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		err = executeSQLFile(db, "init.sql")
		if err != nil {
			log.Fatalf("Error executing SQL file: %v", err)
		}
	})
}

func printDB(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("ERROR querying database", err)
	}
	defer rows.Close()

	fmt.Println("id | username | email | password")
	fmt.Println("---------------------------------")
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.id, &user.username, &user.email, &user.pw); err != nil {
			fmt.Println("Error scanning row: %v", err)
		}
		fmt.Println(user.id, user.username, user.email, user.pw)
	}
}

func AddUser(username string, email string, password string) error {
	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	return err
}

func CheckUser(username string, password string) error {
	return nil
}

func executeSQLFile(db *sql.DB, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var query strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "--") { // Skip comments
			continue
		}
		query.WriteString(line)
		if strings.HasSuffix(line, ";") { // End of SQL statement
			_, err := db.Exec(query.String())
			if err != nil {
				return err
			}
			query.Reset() // Reset query buffer for the next statement
		}
	}
	fmt.Println("SQL file executed successfully")

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
