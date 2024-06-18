package database_test

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	internalDatabase "github.com/matcha-devs/matcha/internal/database"
)

// Use this to run SQL queries directly to help test our internal database
var probe *sql.DB

func setupDBAndOpenSubject(t *testing.T) *internalDatabase.MySQLDatabase {
	// Get the password from the environment, currently using the admin password this may change in the future!!!
	password := os.Getenv("MYSQL_PASSWORD")
	// Connect to MySQL root
	var err error
	probe, err = sql.Open("mysql", "root:"+password+"@tcp(localhost:3306)/")
	if err != nil {
		t.Fatal("Failed to open test database connection -", err)
	}
	if err := probe.Ping(); err != nil {
		t.Fatal("Failed to ping test database connection -", err)
	}

	// Ensure a clean test database by recreating any previous schemas
	_, err = probe.Exec("DROP DATABASE IF EXISTS test_db")
	if err != nil {
		t.Fatal("Failed to drop test database:", err)
	}
	_, err = probe.Exec("CREATE DATABASE test_db")
	if err != nil {
		t.Fatal("Failed to create test database:", err)
	}
	if err := probe.Close(); err != nil {
		log.Fatal("Error closing initial database connection:", err)
	}

	// Connect to the newly created test database
	probe, err = sql.Open("mysql", "root:"+password+"@tcp(localhost:3306)/test_db")
	if err != nil {
		t.Fatal("Failed to open connection to test database:", err)
	}
	if err := probe.Ping(); err != nil {
		t.Fatal("Failed to ping test database connection:", err)
	}

	// Open the custom testing database with initial SQL setup
	subject := internalDatabase.New("test_db", "root", password)
	if subject == nil {
		t.Fatal("Failed to open subject database")
	}
	if err := subject.Open(); err != nil {
		t.Fatal("Failed to open subject database:", err)
	}

	return subject
}

func TestOpenAndClose(t *testing.T) {
	subject := setupDBAndOpenSubject(t)
	if subject == nil {
		t.Error("Database connection is nil")
		return
	}
	// TODO(@Alishah634): Test that all the tables are generated.
	err := subject.Close()
	if err != nil {
		t.Error("Failed to close database:", err)
	}
}

func TestAddUser(t *testing.T) {
	subject := setupDBAndOpenSubject(t)
	defer func() {
		err := subject.Close()
		if err != nil {
			t.Error("Failed to close database:", err)
		}
	}()

	err := subject.AddUser("test_user", "test_user@example.com", "test_pass")
	if err != nil {
		t.Fatal("Failed to add user:", err)
	}

	// Verify the user was added
	var id int
	err = probe.QueryRow("SELECT id FROM users WHERE username = ?", "test_user").Scan(&id)
	if err != nil {
		t.Fatal("Failed to find added user:", err)
	}
	if id == 0 {
		t.Error("Added user ID is zero")
	}
}

func TestAuthenticateLogin(t *testing.T) {
	subject := setupDBAndOpenSubject(t)

	defer func() {
		err := subject.Close()
		if err != nil {
			t.Error("Failed to close database:", err)
		}
	}()

	err := subject.AddUser("test_user", "test_user@example.com", "test_pass")
	if err != nil {
		t.Fatal("Failed to add user:", err)
	}

	t.Run("Valid_Login", func(t *testing.T) {
		// Authenticate valid login
		err = subject.AuthenticateLogin("test_user", "test_pass")
		if err != nil {
			t.Error("Valid login failed:", err)
		}
	})

	t.Run("Invalid_Password", func(t *testing.T) {
		// Authenticate invalid login
		err = subject.AuthenticateLogin("test_user", "wrong_pass")
		if err == nil {
			t.Error("Invalid login did not fail")
		}
	})

	t.Run("Invalid_Username", func(t *testing.T) {
		// Authenticate invalid login
		err = subject.AuthenticateLogin("wrong_user", "test_pass")
		if err == nil {
			t.Error("Invalid login did not fail")
		}
	})

	t.Run("Invalid_Username_and_Password", func(t *testing.T) {
		// Authenticate invalid login
		err = subject.AuthenticateLogin("wrong_user", "wrong_pass")
		if err == nil {
			t.Error("Invalid login did not fail")
		}
	})
}

func TestDeleteUser(t *testing.T) {
	subject := setupDBAndOpenSubject(t)
	defer func() {
		err := subject.Close()
		if err != nil {
			t.Error("Failed to close database:", err)
		}
	}()

	err := subject.AddUser("delete_user", "delete_user@example.com", "delete_pass")
	if err != nil {
		t.Fatal("Failed to add user:", err)
	}

	// Verify the user was added
	var id int
	err = probe.QueryRow("SELECT id FROM users WHERE username = ?", "delete_user").Scan(&id)
	if err != nil {
		t.Fatal("Failed to find added user:", err)
	}
	if id == 0 {
		t.Error("Added user ID is zero")
	}

	// Delete the user
	err = subject.DeleteUser(id)
	if err != nil {
		t.Fatal("Failed to delete user:", err)
	}

	// Verify the user was deleted
	err = probe.QueryRow("SELECT id FROM users WHERE id = ?", id).Scan(&id)
	if err == nil {
		t.Error("Deleted user still exists -", err)
	} else if !errors.Is(err, sql.ErrNoRows) {
		log.Println("Probe failed to verify user -", err)
	}
}

func TestGetUserID(t *testing.T) {
	subject := setupDBAndOpenSubject(t)
	defer func() {
		err := subject.Close()
		if err != nil {
			t.Error("Failed to close database:", err)
		}
	}()

	t.Run(
		"Valid_Users", func(t *testing.T) {
			err := subject.AddUser("user_id_user", "user_id_user@example.com", "user_id_pass")
			if err != nil {
				t.Fatal("Failed to add user:", err)
			}

			// Get user ID by username
			id, err := subject.GetUserID("username", "user_id_user")
			if errors.Is(err, sql.ErrNoRows) {
				t.Error("Failed to get user ID by username. No user found.")
			}
			if id != 1 {
				t.Error("Failed to get user ID by username")
			}
			// Get user ID by email
			id, err = subject.GetUserID("email", "user_id_user@example.com")
			if errors.Is(err, sql.ErrNoRows) {
				t.Error("Failed to get user ID by email. No user found.")
			}
			if id != 1 {
				t.Error("Failed to get user ID by email")
			}
		},
	)

	t.Run(
		"Multiple_Users", func(t *testing.T) {
			log.Println("Testing for multiple users")
			err := subject.AddUser("user2_id_user2", "user2_id_user2@example.com", "user2_id2_pass")
			if err != nil {
				t.Fatal("Failed to add user:", err)
			}

			// Get user ID by username
			id, err := subject.GetUserID("username", "user2_id_user2")
			if errors.Is(err, sql.ErrNoRows) {
				t.Error("Failed to get user ID by username. No user found.")
			}
			if id != 2 {
				t.Error("Failed to get user ID by username")
			}
			// Get user ID by email
			id, err = subject.GetUserID("email", "user2_id_user2@example.com")
			if errors.Is(err, sql.ErrNoRows) {
				t.Error("Failed to get user ID by email. No user found.")
			}
			if id != 2 {
				t.Error("Failed to get user ID by email")
			}
		},
	)

	t.Run(
		"NonExistent_User", func(t *testing.T) {
			log.Println("Testing for when the user does not exist")
			// Get user ID by username
			id, err := subject.GetUserID("username", "user3_id_user3")
			if !errors.Is(err, sql.ErrNoRows) {
				t.Error("Expected to not find a user by username, but found one")
			}
			if id > 0 {
				t.Error("Expected to not find a user by username, but stored ID")
			}
			// Get user ID by email
			id, err = subject.GetUserID("email", "user3_id_user3@example.com")
			if !errors.Is(err, sql.ErrNoRows) {
				t.Error("Expected to not find a user by email, but found one")
			}
			if id > 0 {
				t.Error("Expected to not find a user by email, but stored ID")
			}
		},
	)
}

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory -", err)
	}
	if !strings.HasSuffix(wd, "/matcha") {
		if err := os.Chdir("../.."); err != nil {
			log.Fatal("Failed to change working directory -", err)
		}
	}
	m.Run()
}
