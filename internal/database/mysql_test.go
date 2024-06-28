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
	"golang.org/x/crypto/bcrypt"

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
	if err = probe.Ping(); err != nil {
		if err = probe.Close(); err != nil {
			t.Fatal("Failed to close test database connection -", err)
		}
		t.Fatal("Failed to ping test database connection -", err)
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

func TestConnections(t *testing.T) {
	subject, probe := setup(t)
	if err := subject.underlyingDB.Ping(); err != nil {
		t.Fatal("Failed to open subject db properly -", err)
	}
	if err := probe.Ping(); err != nil {
		t.Fatal("Failed to open probe db properly -", err)
	}
	teardown(t, subject, probe)
	if err := subject.underlyingDB.Ping(); err == nil {
		t.Fatal("Failed to close subject db properly -", err)
	}
	if err := probe.Ping(); err == nil {
		t.Fatal("Failed to close probe db properly -", err)
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

	happyUser := "testUser"
	happyPass := "testPass"
	if err := subject.AddUser(happyUser, "test_user@example.com", happyPass); err != nil {
		t.Fatal("Failed to add", happyUser, "-", err)
	}

	testCases := []struct {
		name       string
		username   string
		password   string
		expectedID int
		happyPath  bool
	}{
		{name: "valid_login", username: happyUser, password: happyPass, expectedID: 1, happyPath: true},
		{name: "bad_user", username: "im a mistake", password: happyPass, expectedID: 0, happyPath: false},
		{name: "bad_pass", username: happyUser, password: "im a mistake", expectedID: 0, happyPath: false},
		{name: "bad_user_and_pass", username: "we're both", password: "mistakes", expectedID: 0, happyPath: false},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				id, err := subject.AuthenticateLogin(testCase.username, testCase.password)
				if (err == nil) != testCase.happyPath {
					mood := "sad"
					if testCase.happyPath {
						mood = "happy"
					}
					t.Error("Was not expecting a", err, "error in this", mood, "test")
				}
				if id != testCase.expectedID {
					t.Error("got id", id, "expected", testCase.expectedID)
				}
			},
		)
	}
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

func TestGetUser(t *testing.T) {
	subject, probe := setup(t)
	defer teardown(t, subject, probe)

	// Add a test user to the database:
	if err := subject.AddUser("test_user", "test_user@example.com","test_pass"); err != nil {
		t.Fatal("Failed to add user -", err)
	}

	// Retrieve the user from the database:
	user := subject.GetUser(1);
	if user == nil {
		t.Fatal("Expected to find user with ID 1, but got nil")
	}

	if user.Username != "test_user" {
		t.Errorf("Expected username to be 'test_user', but got %s", user.Username)
	}
	if user.Email != "test_user@examplecom" {
		t.Errorf("Expected email to be 'test_user@examplecom', but got %s", user.Email)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("test_pass")); err != nil {
		t.Errorf("Password does not match: %v", err)
	}
	if user.CreatedOn.IsZero() {
		t.Errorf("Expected created_on to be set, but got zero value")
	}

	if !user.IsValid() {
		t.Errorf("Expected valid user, but got invalid user: %v", user)
	}

	// Try to retrieve a non-existing user
	nonExistentUser := subject.GetUser(999)
	if nonExistentUser != nil {
		t.Errorf("Expected no user with ID 999, but got: %v", nonExistentUser)
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
			if id != 0 {
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
