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
// var isFirstdb = true

var once sync.Once
var db *sql.DB

// InitDB returns a singleton database instance
// InitDB initializes the database connection
func InitDB() {
    once.Do(func() {
        var err error
        password := os.Getenv("MYSQL_PASSWORD")  // Get the database password from environment variables
        rootDsn := "root:" + password + "@tcp(127.0.0.1:3306)/"

        // Connect to MySQL without specifying a database
        db, err = sql.Open("mysql", rootDsn)
        if err != nil {
            log.Fatalf("Error opening database: %v", err)
        }

        // Ensure connection to MySQL is available
        if err = db.Ping(); err != nil {
            log.Fatalf("Error connecting to MySQL: %v", err)
        }

        // Create the userdb if it does not exist
        _, err = db.Exec("CREATE DATABASE IF NOT EXISTS userdb")
        if err != nil {
            log.Fatalf("Error creating database 'userdb': %v", err)
        }

        // Close the initial connection and reconnect using the specific database
        db.Close()
        userDbDsn := "root:" + password + "@tcp(127.0.0.1:3306)/userdb"
        db, err = sql.Open("mysql", userDbDsn)
        if err != nil {
            log.Fatalf("Error opening userdb database: %v", err)
        }

        if err = db.Ping(); err != nil {
            log.Fatalf("Error connecting to userdb: %v", err)
        }

        // Execute SQL file to configure the userdb
        err = executeSQLFile("init.sql")
        if err != nil {
            log.Fatalf("Error executing SQL file 'init.sql': %v", err)
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
		if err := rows.Scan(&user.id, &user.username, &user.email, &user.password); err != nil {
			fmt.Println("Error scanning row: %v", err)
		}
		fmt.Println(user.id, user.username, user.email, user.password)
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
		return fmt.Errorf("invalid password")
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
