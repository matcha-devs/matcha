package main

func Tester() error {
	printUsersTable()
	err := AddUser("carlos", "cotera_junior@gmail.com", "MEXICAN")
	return err
}
