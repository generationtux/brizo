package app

import (
	"net/http"
	"os"
)

func listenAndServe(addr string, handler http.Handler) error {
	return nil
}

// getAddress gets the port the app should listen on
func getAddress() string {
	port := os.Getenv("BRIZO_PORT")
	if port == "" {
		port = "8080"
	}

	return ":" + port
}
