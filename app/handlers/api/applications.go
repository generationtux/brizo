package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/resources"
	"github.com/go-zoo/bone"
)

// ApplicationIndex provides a listing of all Applications
func ApplicationIndex(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}
	apps, err := resources.AllApplications(db)
	if err != nil {
		log.Printf("Error when retrieving applications: '%s'\n", err)
		http.Error(w, "there was an error when retrieving applications", http.StatusInternalServerError)
		return
	}
	for i := range apps {
		if len(apps[i].Pods) == 0 {
			apps[i].Pods = make([]resources.Pod, 0)
		}
	}
	json.NewEncoder(w).Encode(apps)
}

// ApplicationShow provides an Application
func ApplicationShow(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	id := bone.GetValue(r, "uuid")
	app, err := resources.GetApplication(db, id, resources.GetApplicationPods)
	if err != nil {
		log.Printf("Error when retrieving application: '%s'\n", err)
		http.Error(w, "there was an error when retrieving application", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(app)
}

// ApplicationCreate creates a new Application
func ApplicationCreate(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	var createForm struct {
		Name string
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&createForm)
	defer r.Body.Close()
	if err != nil {
		log.Printf("decoding error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to parse the form", http.StatusInternalServerError)
		return
	}

	app := resources.Application{
		Name: createForm.Name,
	}
	_, err = resources.CreateApplication(db, &app)
	// @todo handle failed save w/out error?
	if err != nil {
		log.Printf("Error when retrieving application: '%s'\n", err)
		http.Error(w, "there was an error when retrieving application", http.StatusInternalServerError)
		return
	}

	// @todo return some sort of content?
	w.WriteHeader(http.StatusCreated)
	return
}
