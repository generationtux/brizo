package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/generationtux/brizo/config"
	"github.com/generationtux/brizo/kube"
	"github.com/urfave/cli"
)

// ConfigureCLI sets up the CLI app
func ConfigureCLI() *cli.App {
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
	router := ConfigureRoutes()
	address := getAddress()

	log.Printf("==> Brizo is starting on %s\n", address)
	http.ListenAndServe(address, router)
	return nil
}

// getAddress gets the port the app should listen on
func getAddress() string {
	port := os.Getenv("BRIZO_PORT")
	if port == "" {
		port = "8080"
	}

	return fmt.Sprintf(":%s", port)
}

// runFlags returns the configuration flags for the run command
func runFlags() []cli.Flag {
	dataFlags := []cli.Flag{
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
			Destination: &config.App.OAuthGithubClientId,
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
