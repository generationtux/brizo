package app

import "net/http"

// CLIRunner runs the CLI application
type CLIRunner func([]string) error

// HTTPServer listens for and handles HTTP requests
type HTTPServer func(string, http.Handler) error

// Initializer verifies the app is ready to be run
type Initializer func() error
