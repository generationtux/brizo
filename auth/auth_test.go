package auth

import (
	"regexp"
	"testing"

	"github.com/generationtux/brizo/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestAuthStateStringGeneration(t *testing.T) {
	stateString := GetOAuthStateString()
	matches, _ := regexp.Match("[a-zA-Z0-9]", []byte(stateString))

	assert.Len(t, stateString, 64)
	assert.True(t, true, matches)
}

func TestHydratesConfig(t *testing.T) {
	config.App.OAuthGithubClientId = "exampleId"
	config.App.OAuthGithubClientSecret = "exampleSecret"
	oauthConf := oauth2.Config{}

	assert.Equal(t, "", oauthConf.ClientID)
	assert.Equal(t, "", oauthConf.ClientSecret)
	HydrateOAuthConfig(&oauthConf)
	assert.Equal(t, "exampleId", oauthConf.ClientID)
	assert.Equal(t, "exampleSecret", oauthConf.ClientSecret)
}
