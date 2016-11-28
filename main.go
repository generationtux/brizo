package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CloudyKit/jet"
	"github.com/go-zoo/bone"
)

var views = jet.NewHTMLSet("./views")

func main() {
	router := bone.New()

	router.HandleFunc("/app", uiHandler)
	router.HandleFunc("/app/*", uiHandler)

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
