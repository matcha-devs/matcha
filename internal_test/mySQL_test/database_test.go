package mySQL_test

import (
	"database/sql"
	"log"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/matcha-devs/matcha/internal/mySQL"
)

// This a variable to hold the test database connection,
// and run all functions relating to the database:
var test_db_connection *sql.DB

// Mock environment variable for testing
func init() {
	// Set the MySQL password for testing
	fmt.Println("Setting up test environment")

	// test_password := os.Getenv("TEST_MYSQL_PASSWORD")
	// fmt.Println("Test password:", test_password)
	// password := os.Getenv("MYSQL_PASSWORD")
	// fmt.Println(("Password:"), password)
}

// Helper function to create a test database connection and ensure it's clean
func setupTestDB(t *testing.T) *mySQL.Database {
	password := os.Getenv("TEST_MYSQL_PASSWORD")
	// Connect to MySQL without specifying matchaDB
	var err error
	test_db_connection, err = sql.Open("mysql", "root:"+password+"@tcp(localhost:3306)/")
	if err != nil {
		t.Fatal("Failed to open test database connection:", err)
	}
	if err := test_db_connection.Ping(); err != nil {
		t.Fatal("Failed to ping test database connection:", err)
	}

	// Open the custom testing database
	test_db := mySQL.Open("test_db", "../../internal/mySQL/queries/")
	if test_db == nil {
		t.Fatal("Failed to open test database")
	}
	
	// Create the matcha_db if it does not exist
	_, err = test_db_connection.Exec("CREATE DATABASE IF NOT EXISTS " + "test_db")
	if err != nil {
		log.Fatal("Error opening Database-", err)
	}
	if err := test_db_connection.Close(); err != nil {
		log.Fatal("Error closing Database-", err)
	}

	// // Clean the database before running tests
	// _, err = test_db_connection.Exec("DELETE FROM users")
	// if err != nil {
	// 	t.Fatal("Failed to clean test database:", err)
	// }

	return test_db
}

// TestOpenAndClose tests the Open and Close functions
func TestOpenAndClose(t *testing.T) {
	test_db := setupTestDB(t)
	if test_db == nil {
		t.Error("Database connection is nil")
		return
	}
	err := test_db.Close()
	if err != nil {
		t.Error("Failed to close database:", err)
	}
}

// TestAddUserAndAuthenticate tests adding a user and authenticating login
func TestAddUserAndAuthenticate(t *testing.T) {
	test_db := setupTestDB(t)
	defer test_db.Close()

	// Add a user
	err := test_db.AddUser("testuser", "testuser@example.com", "testpass")
	if err != nil {
		t.Fatal("Failed to add user:", err)
	}

	// Verify the user was added
	var id int
	err = test_db_connection.QueryRow("SELECT id FROM users WHERE username = ?", "testuser").Scan(&id)
	if err != nil {
		t.Fatal("Failed to find added user:", err)
	}
	if id == 0 {
		t.Error("Added user ID is zero")
	}

	// Authenticate valid login
	err = test_db.AuthenticateLogin("testuser", "testpass")
	if err != nil {
		t.Error("Valid login failed:", err)
	}

	// Authenticate invalid login
	err = test_db.AuthenticateLogin("testuser", "wrongpass")
	if err == nil {
		t.Error("Invalid login did not fail")
	}
}

// TestDeleteUser tests adding and then deleting a user
func TestDeleteUser(t *testing.T) {
	test_db := setupTestDB(t)
	defer test_db.Close()

	// Add a user
	err := test_db.AddUser("deleteuser", "deleteuser@example.com", "deletepass")
	if err != nil {
		t.Fatal("Failed to add user:", err)
	}

	// // Verify the user was added
	var id int
	// err = test_db_connection.QueryRow("SELECT id FROM users WHERE username = ?", "deleteuser").Scan(&id)
	// if err != nil {
	// 	t.Fatal("Failed to find added user:", err)
	// }
	// if id == 0 {
	// 	t.Error("Added user ID is zero")
	// }

	// Delete the user
	err = test_db.DeleteUser(id)
	if err != nil {
		t.Fatal("Failed to delete user:", err)
	}

	// Verify the user was deleted
	err = test_db_connection.QueryRow("SELECT id FROM users WHERE id = ?", id).Scan(&id)
	if err == nil || err != sql.ErrNoRows {
		t.Error("Deleted user still exists or unexpected error occurred:", err)
	}
}

// TestGetUserID tests retrieving a user ID
func TestGetUserID(t *testing.T) {
	test_db := setupTestDB(t)
	defer test_db.Close()

	// Add a user
	err := test_db.AddUser("useriduser", "useriduser@example.com", "useridpass")
	if err != nil {
		t.Fatal("Failed to add user:", err)
	}

	// Get user ID by username
	id := test_db.GetUserID("username", "useriduser")
	if id == 0 {
		t.Error("Failed to get user ID by username")
	}

	// Get user ID by email
	id = test_db.GetUserID("email", "useriduser@example.com")
	if id == 0 {
		t.Error("Failed to get user ID by email")
	}
}
