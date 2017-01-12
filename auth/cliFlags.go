package auth

import (
	"github.com/generationtux/brizo/config"
	"github.com/urfave/cli"
)

// CLIFlags return the auth configuration CLI flags
func CLIFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "app-url",
			Usage:       "public URL for Brizo app",
			Destination: &config.App.URL,
			EnvVar:      "BRIZO_URL",
		},
		cli.StringFlag{
			Name:        "oauth-github-client-id",
			Usage:       "oauth client id of the github application",
			Destination: &config.App.OAuthGithubClientID,
			EnvVar:      "OAUTH_GITHUB_CLIENT_ID",
		},
		cli.StringFlag{
			Name:        "oauth-github-client-secret",
			Usage:       "oauth client secret of the github application",
			Destination: &config.App.OAuthGithubClientSecret,
			EnvVar:      "OAUTH_GITHUB_CLIENT_SECRET",
		},
	}
}
