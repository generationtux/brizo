package api

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/Machiel/slugify"
	"github.com/generationtux/brizo/app/handlers/jsonutil"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/generationtux/brizo/resources"
	"github.com/go-zoo/bone"
	"github.com/pborman/uuid"
)

// VersionIndex provides a listing of all Versions
func VersionIndex(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		jsonutil.DatabaseConnectError().Render(w)
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
		jsonutil.DatabaseConnectError().Render(w)
		return
	}

	client, err := kube.New()
	if err != nil {
		jsonutil.KubeClientConnectionError().Render(w)
		return
	}

	version, err := resources.GetVersion(db, bone.GetValue(r, "version-uuid"), client, true)
	// @todo this could be a 404, so we need to update this error to handle this
	if err != nil {
		log.Printf("Error when retrieving version: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving version")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(version)
}

// VersionCreate creates a new Version
func VersionCreate(w http.ResponseWriter, r *http.Request) {
	var createForm struct {
		Name       string                `json:"name"`
		Replicas   int                   `json:"replicas"`
		Containers []resources.Container `json:"containers"`
		Volumes    []resources.Volume    `json:"volumes"`
	}
	err := jsonRequest(r, &createForm)
	if err != nil {
		jre := jsonutil.NewJSONResponseError(
			http.StatusBadRequest,
			err.Error())
		jsonutil.RespondJSONError(w, jre)
		return
	}

	validContainerName := regexp.MustCompile(`^([a-z0-9_-]+)$`)
	existingContainerNames := make(map[string]bool)
	for index := 0; index < len(createForm.Containers); index++ {
		name := createForm.Containers[index].Name
		if !validContainerName.MatchString(name) {
			jre := jsonutil.NewJSONResponseError(
				// 422 because Docker will not parse container names outside of this pattern
				http.StatusUnprocessableEntity,
				"Invalid container name (%s), only [a-zA-Z0-9_-] characters are allowed", name)
			jsonutil.RespondJSONError(w, jre)
			return
		}
		if existingContainerNames[name] {
			jre := jsonutil.NewJSONResponseError(
				// 422 because Docker will not parse container names outside of this pattern
				http.StatusBadRequest,
				"Duplicate container name (%s) used, names must be unique per container", name)
			jsonutil.RespondJSONError(w, jre)
			return
		}
		existingContainerNames[name] = true
	}

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		jsonutil.DatabaseConnectError().Render(w)
		return
	}

	environment, err := resources.GetEnvironment(db, bone.GetValue(r, "environment-uuid"))
	if err != nil {
		jre := jsonutil.NewJSONResponseError(
			http.StatusNotFound,
			"Environment not found")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	client, err := kube.New()
	if err != nil {
		jsonutil.KubeClientConnectionError().Render(w)
		return
	}

	version := resources.Version{
		UUID:            uuid.New(),
		Name:            createForm.Name,
		Slug:            slugify.Slugify(createForm.Name),
		Containers:      createForm.Containers,
		Replicas:        createForm.Replicas,
		Volumes:         createForm.Volumes,
		EnvironmentID:   environment.ID,
		EnvironmentUUID: environment.UUID,
		Environment:     environment,
	}

	_, err = resources.CreateVersion(db, client, &version)
	if err != nil {
		log.Printf("Error creating version: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"Unable to create version")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	jsonResponse(w, version, 201)
}
