package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // dialect to use for Gorm
)

// Connect will open a new connection to the database
func Connect() (*gorm.DB, error) {
	return gorm.Open("mysql", getDBAddress())
}

// getDBAddress will put together they mysql connection string
func getDBAddress() string {
	user := os.Getenv("BRIZO_DB_USER")
	pass := os.Getenv("BRIZO_DB_PASSWORD")
	host := os.Getenv("BRIZO_DB_HOST")
	db := os.Getenv("BRIZO_DB_NAME")

	port := os.Getenv("BRIZO_DB_PORT")
	if port == "" {
		port = "3306"
	}

	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, db)
}
