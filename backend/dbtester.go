// Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

package backend

func Tester() {
	PrintUsersTable()
	printOpenidTable()
	AddUser("clo", "cotera_hh@gmail.com", "MEXICAN")
	PrintUsersTable()
}
