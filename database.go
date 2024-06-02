package main

import (
	"log"
	"os"
	"sync"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var once sync.Once
var instance *sql.DB

// GetDB returns a singleton database instance
func GetDB() *sql.DB {
	once.Do(func() {
		var err error
		pswd := os.Getenv("MYSQL_PASSWORD") // Ensure this environment variable is set
		dsn := "root:" + pswd + "@tcp(localhost:3306)/userdb"
		instance, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}
		if err = instance.Ping(); err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
	})
	return instance
}
