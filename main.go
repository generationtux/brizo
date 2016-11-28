package main

import (
	"os"

	"github.com/generationtux/brizo/app"
	"github.com/joho/godotenv"
)

func main() {
	// .env file for local configuration during development (see .env.example)
	godotenv.Load()

	cli := app.ConfigureCLI()
	cli.Run(os.Args)
}
