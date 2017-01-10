package app

import "net/http"

// New creates a new application instance
func New() *Application {
	app := new(Application)
	app.version = "0.1.0" // @todo read version dynamically
	app.server = listenAndServe
	app.init = initializeApp
	app.runner = runCli

	return app
}

// Application top level properties and methods for running Brizo
type Application struct {
	version string
	server  HTTPServer
	init    Initializer
	runner  CLIRunner
}

// Version gets the apps current version
func (app *Application) Version() string {
	return app.version
}

// Server starts the HTTP server for the app
func (app *Application) Server(addr string, handler http.Handler) error {
	return app.server(addr, handler)
}

// Initialize makes sure the app is ready to run
func (app *Application) Initialize() error {
	return app.init()
}

// Run executes the CLI
func (app *Application) Run(args []string) error {
	return app.runner(args)
}
