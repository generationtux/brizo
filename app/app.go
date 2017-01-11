package app

import (
	"net/http"
	"os"

	"github.com/generationtux/brizo/app/routes"
	"github.com/generationtux/brizo/auth"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/generationtux/brizo/resources"
)

// Server is a function that starts an HTTP listener and handles requests using the provided handler
type Server func(string, http.Handler) error

// ChecksHealth interface for health checking components
type ChecksHealth func() error

// RunsMigrations interface for running migrations
type RunsMigrations func(...interface{}) error

// Application top level properties and methods for running Brizo
type Application struct {
	serverListener Server
	serverHandler  http.Handler
	healthChecks   []ChecksHealth
	migrator       RunsMigrations
	shouldMigrate  []interface{}
}

// New creates a new application instance
func New() *Application {
	brizo := new(Application)
	brizo.serverListener = http.ListenAndServe
	brizo.serverHandler = routes.BuildRouter()
	brizo.healthChecks = []ChecksHealth{
		database.Health,
		kube.Health,
	}
	brizo.migrator = database.Migrate
	brizo.shouldMigrate = []interface{}{
		&auth.User{},
		&resources.Application{},
	}

	return brizo
}

// Server starts the HTTP server for the app
func (app *Application) Server() error {
	return app.serverListener(getAddress(), app.serverHandler)
}

// Initialize makes sure the app is ready to run
func (app *Application) Initialize() error {
	var err error
	for _, check := range app.healthChecks {
		if err = check(); err != nil {
			return err
		}
	}

	return app.migrator(app.shouldMigrate...)
}

// getAddress gets the port the app should listen on
// @todo refactor to cli flag
func getAddress() string {
	port := os.Getenv("BRIZO_PORT")
	if port == "" {
		port = "8080"
	}

	return ":" + port
}
