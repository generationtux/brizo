package app

import (
	"net/http"

	"github.com/urfave/cli"
)

// Application top level properties and methods for running Brizo
type Application struct {
	version string
	server  HTTPServer
	init    Initializer
	builder CLIBuilder
}

// New creates a new application instance
func New() *Application {
	brizo := new(Application)
	brizo.version = "0.1.0" // @todo read version dynamically
	brizo.server = listenAndServe
	brizo.init = initializeApp
	brizo.builder = buildCli

	return brizo
}

// Version gets the apps current version
func (app *Application) Version() string {
	return app.version
}

// BuildCLI builds the CLI application instance
func (app *Application) BuildCLI() *cli.App {
	return app.builder(app.version)
}

// Server starts the HTTP server for the app
func (app *Application) Server(addr string, handler http.Handler) error {
	return app.server()
}

// Initialize makes sure the app is ready to run
func (app *Application) Initialize() error {
	return app.init()
}

// Run executes the CLI
func (app *Application) Run(args []string) error {
	c := app.BuildCLI()
	return c.Run(args)
}
