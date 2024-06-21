package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// Get the password from the environment, currently using the admin password this may change in the future!
var password = os.Getenv("MYSQL_PASSWORD")

func setup(t *testing.T) (subject *MySQLDatabase, probe *sql.DB) {
	t.Helper()

	// Connect to MySQL root
	probe, err := sql.Open("mysql", "root:"+password+"@tcp(localhost:3306)/")
	if err != nil {
		t.Fatal("Failed to open test database connection -", err)
	}
	if err := probe.Ping(); err != nil {
		_ = probe.Close()
		t.Fatal("Failed to ping test database connection", err)
	}

	// Ensure a clean test database by recreating any previous schemas
	_, err = probe.Exec("DROP DATABASE IF EXISTS test_db")
	if err != nil {
		t.Fatal("Failed to drop test database -", err)
	}
	_, err = probe.Exec("CREATE DATABASE test_db")
	if err != nil {
		t.Fatal("Failed to create test database -", err)
	}

	// Connect to the newly created test database
	_, err = probe.Exec("USE test_db")
	if err != nil {
		log.Fatalln("Failed to use test_db -", err)
	}

	// Open the internal database implementation to test
	subject = New("test_db", "root", password)
	if err := subject.Open(); err != nil {
		t.Fatal("Failed to open subject database -", err)
	}
	return subject, probe
}

func teardown(t *testing.T, subject *MySQLDatabase, probe *sql.DB) {
	t.Helper()

	err := probe.Close()
	if err != nil {
		t.Fatal("Failed to close probe database connection -", err)
	}
	err = subject.Close()
	if err != nil {
		t.Fatal("Failed to close subject database connection -", err)
	}
}

func TestNew(t *testing.T) {
	subject, probe := setup(t)
	defer teardown(t, subject, probe)

	expectedTables := map[string]map[string]struct{}{
		"users":  {"id": {}, "username": {}, "email": {}, "password": {}, "created_on": {}},
		"openid": {"id": {}, "created_on": {}},
	}

	tables, err := probe.Query("SHOW TABLES FROM test_db")
	if err != nil {
		t.Fatal("Failed to query for all tables -", err)
	}
	for tables.Next() {
		var table string
		if err := tables.Scan(&table); err != nil || table == "" {
			t.Fatal("Failed to scan tables -", err)
		}
		expectedCols, exists := expectedTables[table]
		if !exists {
			t.Fatal("Unexpected table:", table)
		}
		t.Run(
			"verify_columns:"+table, func(t *testing.T) {
				columns, err := probe.Query(fmt.Sprintf("SHOW COLUMNS FROM test_db.%s", table))
				if err != nil {
					t.Fatal("Failed to query for columns -", err)
				}
				for columns.Next() {
					var field, typ, null, key, def, extra sql.NullString
					err := columns.Scan(&field, &typ, &null, &key, &def, &extra)
					if err != nil {
						t.Fatal("Failed to scan column -", err)
					}
					delete(expectedCols, field.String)
				}
				if len(expectedCols) != 0 {
					t.Log("Failed to create columns:")
					for k := range expectedCols {
						t.Log(k)
					}
					t.Fail()
				}
			},
		)
		delete(expectedTables, table)
	}
	if len(expectedTables) != 0 {
		t.Log("Failed to create tables:")
		for table := range expectedTables {
			t.Log(table)
		}
		t.Fail()
	}
	if err := tables.Close(); err != nil {
		t.Fatal("Failed to close rows -", err)
	}
}

func TestOpen(t *testing.T) {
	subject, probe := setup(t)
	defer teardown(t, subject, probe)
}

func TestClose(t *testing.T) {
	subject, probe := setup(t)
	defer teardown(t, subject, probe)
}

func TestAddUser(t *testing.T) {
	subject, probe := setup(t)
	defer teardown(t, subject, probe)

	if err := subject.AddUser("test_user", "test_user@example.com", "test_pass"); err != nil {
		t.Fatal("Failed to add user -", err)
	}
	var id int
	if err := probe.QueryRow("SELECT id FROM users WHERE username = ?", "test_user").Scan(&id); err != nil || id != 1 {
		t.Fatal("Failed to create first user with id 1 -", err)
	}

	if err := subject.AddUser("test_user2", "test_user2@example2.com", "test_pass2"); err != nil {
		t.Fatal("Failed to add user -", err)
	}
	if err := probe.QueryRow("SELECT id FROM users WHERE username = ?", "test_user2").Scan(&id); err != nil || id != 2 {
		t.Fatal("Failed to sequentially create user with id 2 -", err)
	}
}

func TestAuthenticateLogin(t *testing.T) {
	subject, probe := setup(t)
	defer teardown(t, subject, probe)

	err := subject.AddUser("test_user", "test_user@example.com", "test_pass")
	if err != nil {
		t.Fatal("Failed to add user -", err)
	}

	t.Run(
		"Valid_Login", func(t *testing.T) {
			// Authenticate valid login
			id, err := subject.AuthenticateLogin("test_user", "test_pass")
			if err != nil || id != 1 {
				t.Error("Valid login failed:", err)
			}
		},
	)

	t.Run(
		"Invalid_Password", func(t *testing.T) {
			// Authenticate invalid login
			id, err := subject.AuthenticateLogin("test_user", "wrong_pass")
			if err == nil || id != 0 {
				t.Error("Invalid login did not fail with id:", id, "-", err)
			}
		},
	)

	t.Run(
		"Invalid_Username", func(t *testing.T) {
			// Authenticate invalid login
			id, err := subject.AuthenticateLogin("wrong_user", "test_pass")
			if err == nil || id != 0 {
				t.Error("Invalid login did not fail with id:", id, "-", err)
			}
		},
	)

	t.Run(
		"Invalid_Username_and_Password", func(t *testing.T) {
			// Authenticate invalid login
			id, err := subject.AuthenticateLogin("wrong_user", "wrong_pass")
			if err == nil || id != 0 {
				t.Error("Invalid login did not fail with id:", id, "-", err)
			}
		},
	)
}

func TestDeleteUser(t *testing.T) {
	subject, probe := setup(t)
	defer teardown(t, subject, probe)

	err := subject.AddUser("delete_user", "delete_user@example.com", "delete_pass")
	if err != nil {
		t.Fatal("Failed to add user -", err)
	}

	// Verify the user was added
	var id int
	err = probe.QueryRow("SELECT id FROM users WHERE username = ?", "delete_user").Scan(&id)
	if err != nil || id == 0 {
		t.Fatal("Failed to find added user -", err)
	}

	// Delete the user
	err = subject.DeleteUser(id)
	if err != nil {
		t.Fatal("Failed to delete user -", err)
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
	subject, probe := setup(t)
	defer teardown(t, subject, probe)

	t.Run(
		"Valid_Users", func(t *testing.T) {
			if err := subject.AddUser("user_id_user", "user_id_user@example.com", "user_id_pass"); err != nil {
				t.Fatal("Failed to add user -", err)
			}

			// Get user ID by username
			id := subject.GetUserID("username", "user_id_user")
			if id != 1 {
				t.Error("Failed to get user ID by username")
			}

			// Get user ID by email
			id = subject.GetUserID("email", "user_id_user@example.com")
			if id != 1 {
				t.Error("Failed to get user ID by email")
			}
		},
	)

	t.Run(
		"Multiple_Users", func(t *testing.T) {
			log.Println("Testing for multiple users")
			if err := subject.AddUser("user2_id_user2", "user2_id_user2@example.com", "user2_id2_pass"); err != nil {
				t.Fatal("Failed to add user -", err)
			}

			// Get user ID by username
			id := subject.GetUserID("username", "user2_id_user2")
			if id != 2 {
				t.Error("Failed to get user ID by username")
			}

			// Get user ID by email
			id = subject.GetUserID("email", "user2_id_user2@example.com")
			if id != 2 {
				t.Error("Failed to get user ID by email")
			}
		},
	)

	t.Run(
		"NonExistent_User", func(t *testing.T) {
			// Get user ID by username
			id := subject.GetUserID("username", "user3_id_user3")
			if id > 0 {
				t.Error("Expected to not find a user by username, but stored ID")
			}

			// Get user ID by email
			id = subject.GetUserID("email", "user3_id_user3@example.com")
			if id != 0 {
				t.Error("Expected to not find a user by email, but stored ID")
			}
		},
	)
}

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed to get working directory -", err)
	}
	if !strings.HasSuffix(wd, "/matcha") {
		if err := os.Chdir("../.."); err != nil {
			log.Fatal("Failed to change working directory -", err)
		}
	}
	m.Run()
}
