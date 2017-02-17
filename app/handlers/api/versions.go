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

// VersionIndex provides a listing of all Versions
func VersionIndex(w http.ResponseWriter, r *http.Request) {
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

	versions, err := resources.GetVersionsByEnvironmentUUID(db, bone.GetValue(r, "environment-uuid"))
	if err != nil {
		log.Printf("Error when retrieving versions: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to connect to the database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(versions)
}

// VersionShow provides a Version
func VersionShow(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(http.StatusInternalServerError, "unable to connect to database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	version, err := resources.GetVersion(db, bone.GetValue(r, "version-uuid"))
	// @todo this could be a 404, so we need to update this error to handle this
	if err != nil {
		log.Printf("Error when retrieving version: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving version")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	environment, err := resources.GetEnvironment(db, bone.GetValue(r, "environment-uuid"))
	if err != nil {
		log.Printf("Error when retrieving version's environment: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving version's environment")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	if version.EnvironmentID != environment.ID {
		jre := jsonutil.NewJSONResponseError(
			http.StatusNotFound,
			"no versions with id of "+version.UUID+" for this environment")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(version)
}

// VersionCreate creates a new Version
func VersionCreate(w http.ResponseWriter, r *http.Request) {
	var createForm struct {
		Name         string                           `json:"name"`
		Image        string                           `json:"image"`
		Replicas     int                              `json:"replicas"`
		Args         []string                         `json:"args"`
		PullPolicy   string                           `json:"pullPolicy"`
		Ports        []resources.ContainerPort        `json:"ports"`
		VolumeMounts []resources.ContainerVolumeMount `json:"volumeMounts"`
		Volumes      []resources.Volume               `json:"volumes"`
	}
	err := jsonRequest(r, &createForm)
	if err != nil {
		jsonErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		jsonErrorResponse(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	environment, err := resources.GetEnvironment(db, bone.GetValue(r, "environment-uuid"))
	if err != nil {
		jsonErrorResponse(w, "Environment not found", http.StatusNotFound)
		return
	}

	client, err := kube.New()
	if err != nil {
		log.Printf("Kube client error: '%s'\n", err)
		jsonErrorResponse(w, "Kube connection error", http.StatusInternalServerError)
		return
	}

	version := resources.Version{
		Name:          createForm.Name,
		Slug:          slugify.Slugify(createForm.Name),
		Image:         createForm.Image,
		Replicas:      createForm.Replicas,
		Volumes:       createForm.Volumes,
		EnvironmentID: environment.ID,
		Environment:   *environment,
		PullPolicy:    createForm.PullPolicy,
		Args:          createForm.Args,
		Ports:         createForm.Ports,
		VolumeMounts:  createForm.VolumeMounts,
	}

	_, err = resources.CreateVersion(db, client, &version)
	if err != nil {
		log.Printf("Error creating version: '%s'\n", err)
		jsonErrorResponse(w, "Unable to create version", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, version, 200)
}
