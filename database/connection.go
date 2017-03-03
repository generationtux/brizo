package database

import (
	"fmt"

	"github.com/generationtux/brizo/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // dialect to use for Gorm
)

// Connect will open a new connection to the database
func Connect() (*gorm.DB, error) {
	return gorm.Open("mysql", getDBAddress())
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
