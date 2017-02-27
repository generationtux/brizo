package auth

import (
	"github.com/generationtux/brizo/config"
	githuboauth "github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

// GetOAuthStateString provides a randomized string to prevent csrf
func GetOAuthStateString() (oauthStateString string) {
	// @todo handle persistent random state
	return "yoKBV7eUR2cVYjx4fsYs3ZPkLwbMfBHg"
}

// BuildUserFromGithubUser handles the checks needed to build a brizo user from
// possible github user attributes.
func BuildUserFromGithubUser(githubUser *githuboauth.User, email string, token string) User {
	var name string
	if githubUser.Name != nil {
		name = *githubUser.Name
	}
	return User{
		Username:       *githubUser.Login,
		Name:           name,
		Email:          email,
		GithubUsername: *githubUser.Login,
		GithubToken:    token,
	}
}

// CreateNewGithubUser takes oauth response values and creates a new Brizo user
func CreateNewGithubUser(db *gorm.DB, githubUser *githuboauth.User, email string, token string) (User, error) {
	user := BuildUserFromGithubUser(githubUser, email, token)
	err := db.Create(&user).Error

	return user, err
}

// HydrateOAuthConfig is used to set the ClientID & ClientSecret at runtime so
// that we can load them via the cli config values
func HydrateOAuthConfig(oauthConf *oauth2.Config) {
	oauthConf.ClientID = config.App.OAuthGithubClientID
	oauthConf.ClientSecret = config.App.OAuthGithubClientSecret
	oauthConf.RedirectURL = config.App.URL + "/o/auth/callback/github"
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
