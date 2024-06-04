// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

import (
	"database/sql"
)

var id = 4

func GetUsers() map[int]User {
	return users
}

func GetUser(id int) User {
	return users[id]
}

// GetPassword retrieves the password for a given user ID
func GetPassword(db *sql.DB, id int) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM users WHERE id = ?", id).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func AddUser(user User) {
	user.ID = id
	id++
	users[user.ID] = user
}

func RemoveUser(id int) {
	delete(users, id)
}

var users = map[int]User{
	0: {ID: 0, Name: "Ancient One", Email: "ancientone@gmail.com", Password: "pw1"},
	1: {ID: 1, Name: "Alice", Email: "alice@example.com", Password: "pw2"},
	2: {ID: 2, Name: "Bob", Email: "bob@example.com", Password: "pw3"},
	3: {ID: 3, Name: "Charlie", Email: "charlie@example.com", Password: "pw4"},
}
