package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/generationtux/brizo/app"
	"github.com/joho/godotenv"
)

func main() {
	// .env file for local configuration during development (see .env.example)
	godotenv.Load()

	router := app.ConfigureRoutes()
	address := getAddress()

	log.Printf("Brizo is starting on %s\n", address)
	http.ListenAndServe(address, router)
}

func getAddress() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return fmt.Sprintf(":%s", port)
}
