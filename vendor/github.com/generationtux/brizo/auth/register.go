package auth

import (
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
func CreateNewGithubUser(db *gorm.DB, githubUser *githuboauth.User, token string) error {
	user := User{
		Username:       *githubUser.Login,
		Name:           *githubUser.Name,
		Email:          *githubUser.Email,
		GithubUsername: *githubUser.Login,
		GithubToken:    token,
	}

	return db.Create(&user).Error
}

// HydrateOAuthConfig is used to set the ClientID & ClientSecret at runtime so
// that we can load them via the cli config values
func HydrateOAuthConfig(oauthConf *oauth2.Config) {
	oauthConf.ClientID = config.App.OAuthGithubClientID
	oauthConf.ClientSecret = config.App.OAuthGithubClientSecret
}

// isFirstUser will determine if any users have registered yet
func isFirstUser(db *gorm.DB) bool {
	var users []User
	db.Find(&users)

	return len(users) == 0
}

// githubUserAllowed will determine if the provided username from Github is allowed
func githubUserAllowed(db *gorm.DB, username string) bool {
	var users []User
	db.Where("github_username = ?", username).Find(&users)

	return len(users) > 0
}
