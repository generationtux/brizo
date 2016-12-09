package auth

import (
	"github.com/generationtux/brizo/config"
	"github.com/generationtux/brizo/database"
	githuboauth "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GetOAuthStateString provides a randomized string to prevent csrf
func GetOAuthStateString() (oauthStateString string) {
	return generateRandomString(64)
}

// CreateNewGithubUser takes oauth response values and creates a new Brizo user
func CreateNewGithubUser(githubUser *githuboauth.User, token string) {
	user := User{
		Username:       *githubUser.Login,
		Name:           *githubUser.Name,
		Email:          *githubUser.Email,
		GithubUsername: *githubUser.Login,
		GithubToken:    token,
	}
	db, _ := database.Connect()
	defer db.Close()
	db.Create(&user)
}

// HydrateOAuthConfig is used to set the ClientID & ClientSecret at runtime so
// that we can load them via the cli config values
func HydrateOAuthConfig(oauthConf *oauth2.Config) {
	oauthConf.ClientID = config.App.OAuthGithubClientId
	oauthConf.ClientSecret = config.App.OAuthGithubClientSecret
}
