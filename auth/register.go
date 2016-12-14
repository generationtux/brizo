package auth

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/generationtux/brizo/config"
	githuboauth "github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

// GetOAuthStateString provides a randomized string to prevent csrf
func GetOAuthStateString() (oauthStateString string) {
	return generateRandomString(64)
}

// CreateNewGithubUser takes oauth response values and creates a new Brizo user
func CreateNewGithubUser(db *gorm.DB, githubUser *githuboauth.User, token string) (User, error) {
	user := User{
		Username:       *githubUser.Login,
		Name:           *githubUser.Name,
		Email:          *githubUser.Email,
		GithubUsername: *githubUser.Login,
		GithubToken:    token,
	}

	err := db.Create(&user).Error

	return user, err
}

// HydrateOAuthConfig is used to set the ClientID & ClientSecret at runtime so
// that we can load them via the cli config values
func HydrateOAuthConfig(oauthConf *oauth2.Config) {
	oauthConf.ClientID = config.App.OAuthGithubClientID
	oauthConf.ClientSecret = config.App.OAuthGithubClientSecret
}

// IsFirstUser will determine if any users have registered yet
func IsFirstUser(db *gorm.DB) bool {
	var users []User
	db.Find(&users)

	return len(users) == 0
}

// GithubUserAllowed will determine if the provided username from Github is allowed
func GithubUserAllowed(db *gorm.DB, username string) bool {
	var users []User
	db.Where("github_username = ?", username).Find(&users)

	return len(users) > 0
}

// create JWT Token
func CreateToken(user User) (string,error) {
	token := jwt.NewWithClaims(tokenSigningMethod(os.Getenv("JWT_ALGO")), jwt.MapClaims{
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
	})

	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))

	return tokenString, err
}

func tokenSigningMethod(method string) jwt.SigningMethod {
	switch method {
	case "HS256":
		return jwt.SigningMethodHS256
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	default:
		return nil
	}
}
