package config

// appConfiguration data structure for application configuration
type appConfiguration struct {
	MysqlHost               string
	MysqlPort               string
	MysqlUser               string
	MysqlPassword           string
	MysqlDatabase           string
	OAuthGithubClientID     string
	OAuthGithubClientSecret string
}

// App holds the apps configuration
var App = appConfiguration{}
