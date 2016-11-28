package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/generationtux/brizo/database"
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
		},
	}
}

// runApp will initialize the app and start the HTTP listener
func runApp(c *cli.Context) error {
	e := initializeApp()
	if e != nil {
		return e
	}

	router := ConfigureRoutes()
	address := getAddress()

	log.Printf("Brizo is starting on %s\n", address)
	http.ListenAndServe(address, router)
	return nil
}

// initializeApp will validate configuration and migrate the database if necessary
func initializeApp() error {
	e := database.Health()
	if e != nil {
		return e
	}

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
