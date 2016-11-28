package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CloudyKit/jet"
	//"github.com/generationtux/brizo/database"
	"github.com/go-zoo/bone"
	"github.com/joho/godotenv"
)

var views = jet.NewHTMLSet("./views")

func main() {
	// .env file for local configuration during development (see .env.example)
	godotenv.Load()

	router := bone.New()
	router.HandleFunc("/app", uiHandler)
	router.HandleFunc("/app/*", uiHandler)

	// healthz endpoint
	router.HandleFunc("/healthz", healthzHandler)

	address := getAddress()
	fmt.Printf("Brizo is starting on %s", address)
	http.ListenAndServe(address, router)
}

func uiHandler(rw http.ResponseWriter, request *http.Request) {
	view, err := views.GetTemplate("index.jet")

	if err != nil {
		log.Println("Unexpected template err:", err.Error())
	}

	view.Execute(rw, nil, nil)
}

func getAddress() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return fmt.Sprintf(":%s", port)
}

func healthzHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Write([]byte("ok"))
}
