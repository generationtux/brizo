package app

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet"
	"github.com/generationtux/brizo/database"
)

var views = jet.NewHTMLSet("./ui")

// rootHandler redirects to Javascript app
func rootHandler(rw http.ResponseWriter, request *http.Request) {
	http.Redirect(rw, request, "/app", 301)
}

// uiHandler for requests to Javascript app
func uiHandler(rw http.ResponseWriter, request *http.Request) {
	view, err := views.GetTemplate("index.html")

	if err != nil {
		log.Println("Unexpected template err:", err.Error())
	}

	view.Execute(rw, nil, nil)
}

// healthzHandler for health check requests
func healthzHandler(rw http.ResponseWriter, request *http.Request) {
	err := database.Health()

	if err != nil {
		log.Println(err)
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
		return
	}

	rw.Write([]byte("ok"))
}
