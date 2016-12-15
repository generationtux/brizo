package dev

import (
	"log"

	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
)

// seeder is a function that seeds data
type seeder func(*gorm.DB) error

func getSeeders() []seeder {
	return []seeder{}
}

// RunSeeders will run all of the listed database seeders
func RunSeeders() error {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		return err
	}

	for _, s := range getSeeders() {
		err = s(db)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
