package database

import (
	"github.com/CarlosACJ55/matcha/internal/database"
)

func _() {
	err := database.AddUser("clo", "cotera_hh@gmail.com", "MEXICAN")
	if err != nil {
		return
	}
}
