package main

import (
	"database/sql"
	"log"
	"os"
	"sync"
	_ "github.com/go-sql-driver/mysql"
	
	// The following imports to use the First database instance:
	"strings"
	"fmt"
	"bufio"
)

var once sync.Once
var instance *sql.DB

// InitDB returns a singleton database instance
func InitDB() *sql.DB {
	once.Do(func() {
		var err error
		pswd := os.Getenv("MYSQL_PASSWORD") // Ensure this environment variable is set
		dsn := "root:" + pswd + "@tcp(127.0.0.1:3306)/userdb"
		instance, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}
		if err = instance.Ping(); err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		
		isFirstInstance := false 
		// Execute SQL script from file if this is the first time the database is created(or you are running this code):
		if(isFirstInstance){
			err = executeSQLFile(instance, "init.sql")
			if err != nil {
			 log.Fatalf("Error executing SQL file: %v", err)
			}
		}
	})
	return instance
}



func printDB(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var email string
		var password string
		err = rows.Scan(&id, &name, &email, &password)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		log.Printf("User: %d, %s, %s, %s\n", id, name, email, password)
	}
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
