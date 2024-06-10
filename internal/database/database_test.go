package database

import (
	"testing"
)

//to run this file 'go test ./...' in the terminal

// TODO(@Alishah634): Move Tests to Test Modules and activate modules

// TODO(@Alishah634): Fix AuthenticateLogin panic error
func TestAuthenticateLogin(t *testing.T) {
	t.Log("TestAuthenticateLogin Started")

	if AuthenticateLogin("Jane", "pw2") != nil {
		t.Error("Authenticated with incorrect username")
	}
	if AuthenticateLogin("Bob", "pw2") == nil {
		t.Error("Authenticated with incorrect password")
	}
	if AuthenticateLogin("Charlie", "pw4") != nil {
		t.Error("Authentication failed with correct password")
	}
}

// TODO(@Alishah634): implement following API functions
//TestAddUser
//
//TestFindID
//
//TestDeleteUser
