package mySQL_test

import (
	"database/sql"
	_ "fmt"
	"log"
	"os"
	"testing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/matcha-devs/matcha/internal/mySQL"
)

// This a variable to hold the test database connection,
// and run all functions relating to the database:
var test_db_connection *sql.DB

// Mock environment variable for testing
// Helper function to create a test database connection and ensure it's clean
func setupTestDB(t *testing.T) *mySQL.Database {
	// Get the password from the environment, currently using the admin password this may change in the future!!!
	password := os.Getenv("MYSQL_PASSWORD")
	// Connect to MySQL without specifying a database
	var err error
	test_db_connection, err = sql.Open("mysql", "root:"+password+"@tcp(localhost:3306)/")
	if err != nil {
		t.Fatal("Failed to open test database connection:", err)
	}
	if err := test_db_connection.Ping(); err != nil {
		t.Fatal("Failed to ping test database connection:", err)
	}

	// Ensure a clean test database by dropping any previous instances
	_, err = test_db_connection.Exec("DROP DATABASE IF EXISTS test_db")
	if err != nil {
		t.Fatal("Failed to drop test database:", err)
	}
	// Create a new test database
	_, err = test_db_connection.Exec("CREATE DATABASE test_db")
	if err != nil {
		t.Fatal("Failed to create test database:", err)
	}

	// Close the initial connection and reopen with the new database
	if err := test_db_connection.Close(); err != nil {
		log.Fatal("Error closing initial database connection:", err)
	}

	// Connect to the newly created test database
	test_db_connection, err = sql.Open("mysql", "root:"+password+"@tcp(localhost:3306)/test_db")
	if err != nil {
		t.Fatal("Failed to open connection to test database:", err)
	}
	if err := test_db_connection.Ping(); err != nil {
		t.Fatal("Failed to ping test database connection:", err)
	}

	// Open the custom testing database with initial SQL setup
	test_db := mySQL.Open("test_db", "root", password, "../../internal/mySQL/queries/")
	if test_db == nil {
		t.Fatal("Failed to open test database")
	}

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
		return
	}
}

// TestAddUserAndAuthenticate tests adding a user and authenticating login
func TestAddUser(t *testing.T) {
	test_db := setupTestDB(t)
	defer test_db.Close()

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
}

func TestAuthenticateLogin(t *testing.T) {
	test_db := setupTestDB(t)
	defer test_db.Close()

	err := test_db.AddUser("testuser", "testuser@example.com", "testpass")
	if err != nil {
		t.Fatal("Failed to add user:", err)
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

	// Verify the user was added
	var id int
	err = test_db_connection.QueryRow("SELECT id FROM users WHERE username = ?", "deleteuser").Scan(&id)
	if err != nil {
		t.Fatal("Failed to find added user:", err)
	}
	if id == 0 {
		t.Error("Added user ID is zero")
	}

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
