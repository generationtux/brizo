package app

import "github.com/urfave/cli"

// HTTPServer listens for and handles HTTP requests
type HTTPServer func() error

// Initializer verifies the app is ready to be run
type Initializer func() error

// CLIBuilder builds the CLI application with provided version
type CLIBuilder func(string) *cli.App

// FlagProvider gets configured CLI flags
type FlagProvider func() []cli.Flag
