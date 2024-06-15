package mySQL_test

import (
	"testing"

	"github.com/matcha-devs/matcha/internal/mySQL"
)

var db = mySQL.Open("testing_db")

// TODO(@SeoYoung Cho):what do is the panic error you are having?:

// TODO(@Alishah634): Move Tests to Test Modules and activate modules

// TODO(@Alishah634): Fix AuthenticateLogin panic error
func TestAuthenticateLogin(t *testing.T) {
	t.Log("TestAuthenticateLogin Started")

	if db.AuthenticateLogin("Jane", "pw2") != nil {
		t.Error("Authenticated with incorrect username")
	}
	if db.AuthenticateLogin("Bob", "pw2") == nil {
		t.Error("Authenticated with incorrect password")
	}
	if db.AuthenticateLogin("Charlie", "pw4") != nil {
		t.Error("Authentication failed with correct password")
	}
}

// TODO(@Alishah634): implement following API tests
//TestAddUser
//
//TestFindID
//
//TestDeleteUser

func _() {
	err := db.AddUser("clo", "cotera_hh@gmail.com", "MEXICAN")
	if err != nil {
		return
	}
}