package auth

import (
    "os"

	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
    "github.com/dgrijalva/jwt-go"
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

// create JWT Token
func CreateToken(user *User) (string, error) {
    token := jwt.NewWithClaims(tokenSigningMethod(os.Getenv("JWT_ALGO")), jwt.MapClaims{
        "username":         user.Username,
        "name":             user.Name,
        "email":            user.Email,
        "githubusername":   user.GithubUsername,
        "githubtoken":      user.GithubToken,
    })

    tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))

    return tokenString, err
}

func tokenSigningMethod(method string) {
    switch method {
        case "HS256":
            return jwt.SigningMethodHS256
        case "HS384"
            return jwt.SigningMethodHS384
        case "HS512"
            return jwt.SigningMethodHS512
        default:
            return nil
    }
}
