package database

import (
	"fmt"

	"github.com/generationtux/brizo/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // dialect to use for Gorm
)

// Connect will open a new connection to the database. If there is an error when
// attempting to connect to the database, we will close the database connection
// here prior to returning.
func Connect() (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	if db, err = gorm.Open("mysql", getDBAddress()); err != nil {
		db.Close()
	}

	return db, err
}

// getDBAddress will put together the mysql connection string
func getDBAddress() string {
	user := config.App.MysqlUser
	pass := config.App.MysqlPassword
	host := config.App.MysqlHost
	port := config.App.MysqlPort
	db := config.App.MysqlDatabase

	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, db)
}

// Health will check make sure the database is reachable or return an error
func Health() error {
	db, err := Connect()
	defer db.Close()

	return err
}
