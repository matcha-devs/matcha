package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"sync"

	// The following imports to use the First database db:
	"bufio"
	f "fmt"
	"strings"
)

//If you don't have 'userdb' on MySQL, set it to 'true'
var isFirstdb bool = false

var once sync.Once
var db *sql.DB

// InitDB returns a singleton database db
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

		//If you don't have 'userdb' on MySQL, set it to 'true'
		if isFirstdb {
			err := executeSQLFile("init.sql")
			if err != nil {
				log.Fatalf("Error executing SQL file: %v", err)
			}
		}
	})
}

func printDB() {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		f.Println("ERROR querying database", err)
	}
	defer rows.Close()

	f.Println("id | username | email | password")
	f.Println("---------------------------------")
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.id, &user.username, &user.email, &user.pw); err != nil {
			f.Println("Error scanning row: %v", err)
		}
		f.Println(user.id, user.username, user.email, user.pw)
	}
}

func addUser(username string, email string, password string) error {
	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	return err
}

func checkUser(username, password string) (bool, error) {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	f.Println("DB Password:", dbPassword)
	f.Println("err:", err)

	if err != nil {
		return false, err
	}
	return dbPassword == password, nil
}


func executeSQLFile(filepath string) error {
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
	f.Println("SQL file executed successfully")

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
