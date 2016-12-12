package web

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet"
	"github.com/generationtux/brizo/database"
)

var views = jet.NewHTMLSet("./ui")

// RootHandler redirects to Javascript app
func RootHandler(rw http.ResponseWriter, request *http.Request) {
	http.Redirect(rw, request, "/app", 301)
}

// UIHandler for requests to Javascript app
func UIHandler(rw http.ResponseWriter, request *http.Request) {
	view, err := views.GetTemplate("index.html")

	if err != nil {
		log.Println("Unexpected template err:", err.Error())
	}

	view.Execute(rw, nil, nil)
}

// HealthzHandler for health check requests
func HealthzHandler(rw http.ResponseWriter, request *http.Request) {
	err := database.Health()

	if err != nil {
		log.Println(err)
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
		return
	}

	rw.Write([]byte("ok"))
}
