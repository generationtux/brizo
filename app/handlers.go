package app

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet"
	"github.com/generationtux/brizo/database"
)

var views = jet.NewHTMLSet("./views")

// uiHandler for requests to Javascript app
func uiHandler(rw http.ResponseWriter, request *http.Request) {
	view, err := views.GetTemplate("index.jet")

	if err != nil {
		log.Println("Unexpected template err:", err.Error())
	}

	view.Execute(rw, nil, nil)
}

// healthzHandler for health check requests
func healthzHandler(rw http.ResponseWriter, request *http.Request) {
	db, err := database.Connect()
	defer db.Close()

	if err != nil {
		log.Println(err)
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
		return
	}

	rw.Write([]byte("ok"))
}
