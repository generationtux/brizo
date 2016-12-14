package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
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

// CreateJWTToken creates a JWT token for the provided user
func CreateJWTToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwtSigningMethod(), jwt.MapClaims{
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
	})

	tokenString, err := token.SignedString(jwtSecret())

	return tokenString, err
}
