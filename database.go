package main

import (
	"bufio"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"sync"

	"fmt"
)

// If you don't have 'userdb' on MySQL, set it to 'true'
var isFirstdb = false

var once sync.Once
var db *sql.DB

// InitDB returns a singleton database instance
func InitDB() {
	once.Do(func() {
		var err error
		password := os.Getenv("MYSQL_PASSWORD") // Ensure this environment variable is set
		dsn := "root:" + password + "@tcp(127.0.0.1:3306)/userdb"
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}
		if err = db.Ping(); err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		if isFirstdb {
			err = executeSQLFile("init.sql")
			if err != nil {
				log.Fatalf("Error executing SQL file: %v", err)
			}
		}
	})
}

func printUsersTable() {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("ERROR querying database", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

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

func AuthenticateLogin(username, password string) error {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	fmt.Println("DB Password:", dbPassword)
	fmt.Println("err:", err)

	if err != nil {
		return err
	} else if dbPassword != password {
		return fmt.Errorf("Invalid password")
	}
	return nil
}

func executeSQLFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

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
