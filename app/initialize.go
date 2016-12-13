package app

import (
	"log"

	"github.com/generationtux/brizo/auth"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/generationtux/brizo/resources"
)

// initializeApp will validate configuration and migrate the database if necessary
func initializeApp() error {
	_, e := kube.Client()
	if e != nil {
		return e
	}

	e = database.Health()
	if e != nil {
		return e
	}

	return runMigrations()
}

// runMigrations will apply schema changes to the database as necessary
// see http://jinzhu.me/gorm/database.html#migration
func runMigrations() error {
	db, e := database.Connect()
	defer db.Close()
	if e != nil {
		return e
	}

	log.Println("==> Running database migrations...")
	return db.AutoMigrate(
		&auth.User{},
		&resources.Application{},
	).Error
}
