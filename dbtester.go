package main

func Tester() error {
	printUsersTable()
	err := AddUser("notacarlos", "cotera_junior@gmail.com", "MEXICAN")
	return err
}
