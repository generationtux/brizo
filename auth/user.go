package auth

import (
	"github.com/generationtux/brizo/database"
)

// User represents a Brizo user
type User struct {
	database.Model
	Username       string `gorm:"not null;unique_index"`
	Name           string
	Email          string
	GithubUsername string
	GithubToken    string
}

// CreateUser creates a Brizo specific user without forcing any validation on
// the input
func CreateUser(user *User) (bool, error) {
	db, _ := database.Connect()
	defer db.Close()
	result := db.Create(&user)

	return result.RowsAffected == 1, result.Error
}

// UpdateUser updates an existing Brizo user based on his/her username
func UpdateUser(user *User) (bool, error) {
	db, _ := database.Connect()
	defer db.Close()
	result := db.Model(user).Where("username = ?", user.Username).UpdateColumns(user)

	return result.RowsAffected == 1, result.Error
}
