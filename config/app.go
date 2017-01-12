package config

// AppConfiguration data structure for application configuration
type AppConfiguration struct {
	URL                     string
	MysqlHost               string
	MysqlPort               string
	MysqlUser               string
	MysqlPassword           string
	MysqlDatabase           string
	OAuthGithubClientID     string
	OAuthGithubClientSecret string
}

// App holds the apps configuration
var App = AppConfiguration{}
