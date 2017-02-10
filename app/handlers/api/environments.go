package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/generationtux/brizo/app/handlers/jsonutil"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/resources"
	"github.com/go-zoo/bone"
)

// EnvironmentIndex provides a listing of all Environments
func EnvironmentIndex(w http.ResponseWriter, r *http.Request) {
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
	environments, err := resources.AllEnvironments(db)
	if err != nil {
		log.Printf("Error when retrieving environments: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to connect to the database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	for i := range environments {
		if len(environments[i].Versions) == 0 {
			environments[i].Versions = make([]resources.Version, 0)
		}
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(environments)
}

// EnvironmentShow provides an Environment
func EnvironmentShow(w http.ResponseWriter, r *http.Request) {
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

	id := bone.GetValue(r, "uuid")
	environment, err := resources.GetEnvironment(db, id)
	if err != nil {
		log.Printf("Error when retrieving environment: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving environment")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	if len(environment.Versions) == 0 {
		environment.Versions = make([]resources.Version, 0)
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(environment)
}

// EnvironmentCreate creates a new Environment
func EnvironmentCreate(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	var createForm struct {
		Name          string
		ApplicationID string
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&createForm)
	defer r.Body.Close()
	if err != nil {
		log.Printf("decoding error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to parse the form", http.StatusInternalServerError)
		return
	}
	// @todo
	appID, _ := strconv.ParseUint(createForm.ApplicationID, 10, 0)
	environment := resources.Environment{
		Name:          createForm.Name,
		ApplicationID: appID,
	}
	_, err = resources.CreateEnvironment(db, &environment)
	// @todo handle failed save w/out error?
	if err != nil {
		log.Printf("Error when retrieving environment: '%s'\n", err)
		http.Error(w, "there was an error when retrieving environment", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, environment, 200)
}

// EnvironmentEdit edits an Environment
func EnvironmentEdit(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "there was an error when attempting to parse the form", http.StatusInternalServerError)
		return
	}

	id := bone.GetValue(r, "uuid")
	environment, err := resources.GetEnvironment(db, id)
	if err != nil {
		log.Printf("Error when retrieving environment: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving environment")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	environment.Name = editForm.Name

	_, err = resources.UpdateEnvironment(db, environment)
	if err != nil {
		log.Printf("Error when updating environment: '%s'\n", err)
		http.Error(w, "there was an error when updating environment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(environment)
	w.WriteHeader(http.StatusOK)
}
