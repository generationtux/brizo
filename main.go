package main

import (
	"log"

	"github.com/generationtux/brizo/app"
	"github.com/joho/godotenv"
)

func main() {
	// .env file for local configuration during development (see .env.example)
	godotenv.Load()

	brizo := app.New()

	// Intializing Brizo will make sure everything is configured correctly
	// and perform any migration tasks that are required.
	err := brizo.Initialize()
	if err != nil {
		log.Fatalln(err)
	}

	// Run Brizo to handle the current command
}
