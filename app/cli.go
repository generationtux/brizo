package app

import (
	"log"
	"net/http"

	"github.com/generationtux/brizo/app/routes"
	"github.com/generationtux/brizo/config"
	"github.com/generationtux/brizo/kube"
	"github.com/urfave/cli"
)

func runCli(args []string) error {
	return nil
}

// configureCLI sets up the CLI app
func configureCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "Brizo"
	app.Usage = "PAAS powered by Kubernetes."
	app.Version = "0.1.0"
	app.Commands = registerCommands()

	return app
}

// registerCommands sets up the CLI commands
func registerCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "run",
			Usage:  "start Brizo",
			Action: runApp,
			Flags:  runFlags(),
		},
	}
}

// runApp will initialize the app and start the HTTP listener
func runApp(c *cli.Context) error {
	log.Printf("==> Initializing app..")
	e := initializeApp()
	if e != nil {
		log.Println("== error during initilization")
		return e
	}

	log.Printf("==> App is ready")
	router := routes.BuildRouter()
	address := getAddress()

	log.Printf("==> Brizo is starting on %s\n", address)
	http.ListenAndServe(address, router)
	return nil
}

// runFlags returns the configuration flags for the run command
func runFlags() []cli.Flag {
	dataFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "app-url",
			Usage:       "public URL for Brizo app",
			Destination: &config.App.URL,
			EnvVar:      "BRIZO_URL",
		},
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
	kubeFlags := kube.CLIFlags()

	return append(dataFlags, kubeFlags...)
}
