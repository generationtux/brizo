package main

import (
	"fmt"
	"github.com/go-zoo/bone"
	"net/http"
	"os"
)

func main() {
	router := bone.New()

	router.HandleFunc("/app", uiHandler)
	router.HandleFunc("/app/*", uiHandler)

	address := getAddress()
	fmt.Printf("Brizo is starting on %s", address)
	http.ListenAndServe(address, router)
}

func uiHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Write([]byte("Brizo UI"))
}

func getAddress() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return fmt.Sprintf(":%s", port)
}
