package database

import (
	"github.com/generationtux/brizo/config"
	"github.com/urfave/cli"
)

// CLIFlags returns database configuration flags
func CLIFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "mysql-host",
			Usage:       "mysql host address",
			Value:       "localhost",
			Destination: &config.App.MysqlHost,
			EnvVar:      "BRIZO_MYSQL_HOST",
		},
		cli.StringFlag{
			Name:        "mysql-port",
			Usage:       "port mysql is listening on",
			Value:       "3306",
			Destination: &config.App.MysqlPort,
			EnvVar:      "BRIZO_MYSQL_PORT",
		},
		cli.StringFlag{
			Name:        "mysql-user",
			Usage:       "user to connect to mysql",
			Value:       "root",
			Destination: &config.App.MysqlUser,
			EnvVar:      "BRIZO_MYSQL_USER",
		},
		cli.StringFlag{
			Name:        "mysql-password",
			Usage:       "password to connect to mysql",
			Destination: &config.App.MysqlPassword,
			EnvVar:      "BRIZO_MYSQL_PASSWORD",
		},
		cli.StringFlag{
			Name:        "mysql-database",
			Usage:       "name of mysql database to use",
			Value:       "brizo",
			Destination: &config.App.MysqlDatabase,
			EnvVar:      "BRIZO_MYSQL_DATABASE",
		},
	}
}
