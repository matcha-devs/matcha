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
	if err := probe.Ping(); err != nil {
		if err = probe.Close(); err != nil {
			t.Fatal("Failed to close test database connection -", err)
		}
		t.Fatal("Failed to ping test database connection -", err)
	}

	// Ensure a clean test database by recreating any previous schemas
	if _, err := probe.Exec("DROP DATABASE IF EXISTS test_db"); err != nil {
		t.Fatal("Failed to drop test database -", err)
	}
	if _, err := probe.Exec("CREATE DATABASE test_db"); err != nil {
		t.Fatal("Failed to create test database -", err)
	}

	// Connect to the newly created test database
	if _, err := probe.Exec("USE test_db"); err != nil {
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

	testCases := []struct {
		name          string
		username      string
		email         string
		password      string
		expectedID    int
		expectedError bool
	}{
		{"AddFirstUser", "test_user", "test_user@example.com", "test_pass", 1, false},
		{"AddSecondUser", "test_user2", "test_user2@example2.com", "test_pass2", 2, false},
		{"AddDuplicateUsername", "test_user", "unique_email@example.com", "test_pass3", 0, true},
		{"AddDuplicateEmail", "unique_user", "test_user2@example2.com", "test_pass4", 0, true},
		{"AddEmptyUsername", "", "empty_user@example.com", "test_pass5", 0, true},
		{"AddEmptyEmail", "empty_email_user", "", "test_pass6", 0, true},
		{"AddEmptyPassword", "empty_pass_user", "empty_pass_user@example.com", "", 0, true},

		// TODO: The functionality for this test need to be implemented
		// {"AddInvalidEmail", "invalid_email_user", "invalidemail.com", "test_pass7", 0, true},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				err := subject.AddUser(tc.username, tc.email, tc.password)
				var id int
				if tc.expectedError {
					if err == nil {
						t.Fatalf("Expected error but got none for case: %s", tc.name)
					}
				} else {
					if err != nil {
						t.Fatalf("Failed to add user - %v for case: %s", err, tc.name)
					}
					err := probe.QueryRow("SELECT id FROM users WHERE username = ?", tc.username).Scan(&id)
					if err != nil {
						t.Fatalf("Failed to query user id - %v for case: %s", err, tc.name)
					}
					if id != tc.expectedID {
						t.Fatalf("Expected user id %d but got %d for case: %s", tc.expectedID, id, tc.name)
					}
				}
			},
		)
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
	if err := subject.AddUser("test_user", "test_user@example.com", "test_pass"); err != nil {
		t.Fatal("Failed to add user -", err)
	}

	testCases := []struct {
		name       string
		userID     int
		expectUser bool
		username   string
		email      string
		password   string
	}{
		{"GetExistingUser", 1, true, "test_user", "test_user@example.com", "test_pass"},
		{"GetNonExistentUser", 999, false, "", "", ""},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				user := subject.GetUser(tc.userID)
				if tc.expectUser {
					if user == nil {
						t.Fatal("Expected to find user, but got nil")
					}
					if user.Username != tc.username {
						t.Errorf("Expected username to be '%s', but got %s", tc.username, user.Username)
					}
					if user.Email != tc.email {
						t.Errorf("Expected email to be '%s', but got %s", tc.email, user.Email)
					}
					if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tc.password)); err != nil {
						t.Errorf("Password does not match: %v", err)
					}
					if user.CreatedOn.IsZero() {
						t.Errorf("Expected created_on to be set, but got zero value")
					}
					if !user.IsValid() {
						t.Errorf("Expected valid user, but got invalid user: %v", user)
					}
				} else {
					if user != nil {
						t.Errorf("Expected no user, but got: %v", user)
					}
				}
			},
		)
	}
}

func TestGetUserID(t *testing.T) {
	subject, probe := setup(t)
	defer teardown(t, subject, probe)

	t.Log("Adding user: user_id_user")
	if err := subject.AddUser("user_id_user", "user_id_user@example.com", "user_id_pass"); err != nil {
		t.Fatal("Failed to add user -", err)
	}

	t.Log("Adding user: user2_id_user2")
	if err := subject.AddUser("user2_id_user2", "user2_id_user2@example.com", "user2_id2_pass"); err != nil {
		t.Fatal("Failed to add user -", err)
	}

	testCases := []struct {
		name     string
		field    string
		value    string
		expected int
	}{
		{"GetUserIDByUsername1", "username", "user_id_user", 1},
		{"GetUserIDByEmail1", "email", "user_id_user@example.com", 1},
		{"GetUserIDByUsername2", "username", "user2_id_user2", 2},
		{"GetUserIDByEmail2", "email", "user2_id_user2@example.com", 2},
		{"GetNonExistentUserIDByUsername", "username", "nonexistent_user", 0},
		{"GetNonExistentUserIDByEmail", "email", "nonexistent@example.com", 0},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				id := subject.GetUserID(tc.field, tc.value)
				if id != tc.expected {
					t.Errorf("got id %d, expected %d", id, tc.expected)
				}
			},
		)
	}
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
