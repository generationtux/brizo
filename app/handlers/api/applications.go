package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Machiel/slugify"
	"github.com/generationtux/brizo/app/handlers/jsonutil"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/generationtux/brizo/resources"
	"github.com/go-zoo/bone"
)

// ApplicationIndex provides a listing of all Applications
func ApplicationIndex(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to connect to the database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	apps, err := resources.AllApplications(db)
	if err != nil {
		log.Printf("Error when retrieving applications: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to connect to the database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	for i := range apps {
		if len(apps[i].Environments) == 0 {
			apps[i].Environments = make([]resources.Environment, 0)
		}
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

// ApplicationShow provides an Application
func ApplicationShow(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"unable to connect to database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	kubeClient, err := kube.New()
	if err != nil {
		log.Printf("Kube client error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusServiceUnavailable,
			"unable to connect to Kubernetes")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	id := bone.GetValue(r, "uuid")
	app, err := resources.GetApplication(db, kubeClient, id)
	if err != nil {
		log.Printf("Error when retrieving application: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving application")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	if len(app.Environments) == 0 {
		app.Environments = make([]resources.Environment, 0)
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(app)
}

// ApplicationUpdate updates an existing Application
func ApplicationUpdate(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to connect to the database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	var editForm struct {
		Name string
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&editForm)
	defer r.Body.Close()
	if err != nil {
		log.Printf("decoding error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to parse the form")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	kubeClient, err := kube.New()
	if err != nil {
		log.Printf("Kube client error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusServiceUnavailable,
			"unable to connect to Kubernetes")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	id := bone.GetValue(r, "uuid")
	application, err := resources.GetApplication(db, kubeClient, id)
	if err != nil {
		log.Printf("Error when retrieving application: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving application")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	application.Name = editForm.Name
	application.Slug = slugify.Slugify(editForm.Name)

	_, err = resources.UpdateApplication(db, application)
	if err != nil {
		log.Printf("Error when updating application: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when updating application")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(application)
	w.WriteHeader(http.StatusOK)
	return
}

// ApplicationCreate creates a new Application
func ApplicationCreate(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to connect to the database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	var createForm struct {
		Name string
		Slug string
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&createForm)
	defer r.Body.Close()
	if err != nil {
		log.Printf("decoding error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to parse the form")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	app := resources.Application{
		Name: createForm.Name,
		Slug: slugify.Slugify(createForm.Name),
	}
	_, err = resources.CreateApplication(db, &app)
	// @todo handle failed save w/out error?
	if err != nil {
		log.Printf("Error when retrieving application: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving application")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	jsonResponse(w, app, 201)
}
