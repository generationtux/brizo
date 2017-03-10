package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/generationtux/brizo/app/handlers/jsonutil"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
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

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(environment)
}

// EnvironmentCreate creates a new Environment
func EnvironmentCreate(w http.ResponseWriter, r *http.Request) {
	client, err := kube.New()
	if err != nil {
		log.Printf("Kube client error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusServiceUnavailable,
			"unable to reach Kubernetes")
		jsonutil.RespondJSONError(w, jre)
		return
	}

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
		Name          string
		ApplicationID string
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
	// @todo
	appID, _ := strconv.ParseUint(createForm.ApplicationID, 10, 0)
	environment := resources.Environment{
		Name:          createForm.Name,
		ApplicationID: appID,
	}

	app, err := resources.GetApplicationByID(db, createForm.ApplicationID)
	if err != nil {
		log.Printf("Error when retrieving application: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving application")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	_, err = resources.CreateEnvironment(db, client, &environment, app)
	// @todo handle failed save w/out error?
	if err != nil {
		log.Printf("Error when retrieving environment: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving environment")
		jsonutil.RespondJSONError(w, jre)
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
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to parse the form")
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

	environment.Name = editForm.Name

	_, err = resources.UpdateEnvironment(db, environment)
	if err != nil {
		log.Printf("Error when updating environment: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when updating environment")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(environment)
	w.WriteHeader(http.StatusOK)
}

// EnvironmentDeploy deploys a Version to the environment
func EnvironmentDeploy(w http.ResponseWriter, r *http.Request) {
	var deployForm struct {
		// id of the Version we want to deploy
		VersionUUID string `json:"version_uuid"`
	}
	err := jsonRequest(r, &deployForm)
	if err != nil {
		jsonErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		jsonErrorResponse(w, "Unable to connect to database", http.StatusInternalServerError)
		return
	}

	deployToEnvironment, err := resources.GetEnvironment(db, bone.GetValue(r, "environment-uuid"))
	if err != nil {
		log.Println(err)
		jsonErrorResponse(w, "Unable to get environment", http.StatusNotFound)
		return
	}

	client, err := kube.New()
	if err != nil {
		log.Printf("Kube client error: '%s'\n", err)
		jsonErrorResponse(w, "Unable to connect to Kubernetes", http.StatusInternalServerError)
		return
	}

	version, err := resources.GetVersion(db, deployForm.VersionUUID, client, true)
	if err != nil {
		log.Println(err)
		jsonErrorResponse(w, "Unable to get version", http.StatusNotFound)
		return
	}

	// modify the version for the new environment
	version.Environment = deployToEnvironment
	version.EnvironmentID = deployToEnvironment.ID

	_, err = resources.DeployVersion(client, version)
	if err != nil {
		log.Printf("Error deploying version: '%s'\n", err)
		jsonErrorResponse(w, "Unable to deploy version", http.StatusInternalServerError)
		return
	}

	_, err = resources.UpdateVersion(db, version)
	if err != nil {
		log.Printf("Error updating version after deployment: '%s'\n", err)
		jsonErrorResponse(w, "Unable to update version after deployment", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": "Version deployed."}, 200)
}
