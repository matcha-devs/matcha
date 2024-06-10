package database

import (
	"github.com/matcha-devs/matcha/internal/database"
)

func _() {
	err := database.AddUser("clo", "cotera_hh@gmail.com", "MEXICAN")
	if err != nil {
		return
	}
}
