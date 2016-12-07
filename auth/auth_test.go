package auth

import (
	"testing"

	"github.com/generationtux/brizo/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

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
