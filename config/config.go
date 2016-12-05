package config

// configuration data structure for configuration
type configuration struct {
	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPassword string
	MysqlDatabase string
}

// App holds the apps configuration
var App = configuration{}
