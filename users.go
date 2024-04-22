// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package main

func GetUsers() map[int]User {
	return users
}

func GetUser(id int) User {
	return users[id]
}

func AddUser(user User) {
	users[user.ID] = user
}

func RemoveUser(id int) {
	delete(users, id)
}

var users = map[int]User{
	0: {ID: 0, Name: "Ancient One", Email: "ancientone@gmail.com"},
	1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
	2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
	3: {ID: 3, Name: "Charlie", Email: "charlie@example.com"},
}
