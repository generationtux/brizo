package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
)

// User represents a Brizo user
type User struct {
	database.Model
	Username       string `gorm:"not null;unique_index" json:"username"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	GithubUsername string `json:"github_username"`
	GithubToken    string `json:"github_token"`
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
	claims := jwt.MapClaims{
		"exp":      time.Now().Add(time.Hour * 24).UTC().Unix(),
		"iat":      time.Now().UTC().Unix(),
		"iss":      "brizo",
		"sub":      user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"username": user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	secret := []byte(jwtSecret())
	tokenString, err := token.SignedString(secret)

	return tokenString, err
}

// GetInvitedUsers will provide a slice of all of the Users who have been
// invited to Brizo, but have yet to authenticate with Github.
func GetInvitedUsers(db *gorm.DB) ([]User, error) {
	var invitees []User
	result := db.Where("username is not null").Where("github_token = \"\"").Find(&invitees)

	return invitees, result.Error
}
