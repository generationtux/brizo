package auth

import (
    "os"

	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
    jwt "github.com/dgrijalva/jwt-go"
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
func CreateUser(db *gorm.DB, user *User) (bool, error) {
	result := db.Create(&user)

	return result.RowsAffected == 1, result.Error
}

// UpdateUser updates an existing Brizo user based on his/her username
func UpdateUser(db *gorm.DB, user *User) (bool, error) {
	result := db.Model(user).Where("username = ?", user.Username).UpdateColumns(user)

	return result.RowsAffected == 1, result.Error
}

func CreateToken(user *User) (string, error) {
    token := jew.New()
}

func tokenSigningMethod(method string) {
    /*
    HS256
    HS384
    HS512
     */
}
