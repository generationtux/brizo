package auth

import (
	"math/rand"
	"time"

	"golang.org/x/oauth2"

	"github.com/generationtux/brizo/config"
)

func hydrateOAuthConfig(oauthConf *oauth2.Config) {
	oauthConf.ClientID = config.App.OAuthGithubClientId
	oauthConf.ClientSecret = config.App.OAuthGithubClientSecret
}

func generateRandomString(length int) string {
	const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = characters[r.Intn(len(characters))]
	}

	return string(bytes)
}
