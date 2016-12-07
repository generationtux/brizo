package auth

import (
	"github.com/generationtux/brizo/config"
	"github.com/generationtux/brizo/database"
	githuboauth "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetOAuthStateString() (oauthStateString string) {
	return generateRandomString(64)
}

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

func HydrateOAuthConfig(oauthConf *oauth2.Config) {
	oauthConf.ClientID = config.App.OAuthGithubClientId
	oauthConf.ClientSecret = config.App.OAuthGithubClientSecret
}
