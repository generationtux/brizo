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
	"github.com/mholt/binding"
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

// VersionUpdate updates an existing Version
func VersionUpdate(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	var editForm struct {
		Name     string
		Image    string
		Replicas int
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&editForm)
	defer r.Body.Close()
	if err != nil {
		log.Printf("decoding error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to parse the form", http.StatusInternalServerError)
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

	version.Name = editForm.Name
	version.Slug = slugify.Slugify(editForm.Name)
	version.Image = editForm.Image
	version.Replicas = editForm.Replicas

	_, err = resources.UpdateVersion(db, version)
	if err != nil {
		log.Printf("Error when updating version: '%s'\n", err)
		http.Error(w, "there was an error when updating version", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(version)
	w.WriteHeader(http.StatusOK)
	return
}

type versionCreateJSON struct {
	Name          string `json:"name"`
	Image         string `json:"image"`
	Replicas      int    `json:"replicas"`
	Volumes       []resources.Volume
	EnvironmentID string `json:"environment_id"`
}

func (vc *versionCreateJSON) FieldMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&vc.Name:          "name",
		&vc.Image:         "image",
		&vc.Replicas:      "replicas",
		&vc.EnvironmentID: "environment_id",
	}
}

// VersionCreate creates a new Version
func VersionCreate(w http.ResponseWriter, r *http.Request) {
	createVersionJSON := new(versionCreateJSON)
	errs := binding.Bind(r, createVersionJSON)
	if errs.Handle(w) {
		return
	}

	client, err := kube.New()
	if err != nil {
		log.Printf("Kube client error: '%s'\n", err)
		http.Error(w, "unable to reach Kubernetes", http.StatusInternalServerError)
		return
	}

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	environment, err := resources.GetEnvironment(db, bone.GetValue(r, "environment-uuid"))
	if err != nil {
		log.Printf("Error when retrieving environment: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving environment")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	version := resources.Version{
		Name:          createVersionJSON.Name,
		Slug:          slugify.Slugify(createVersionJSON.Name),
		Image:         createVersionJSON.Image,
		Replicas:      createVersionJSON.Replicas,
		EnvironmentID: environment.ID,
		Environment:   *environment,
	}
	_, err = resources.CreateVersion(db, client, &version, createVersionJSON.Volumes)
	// @todo handle failed save w/out error?
	if err != nil {
		log.Printf("Error when retrieving version: '%s'\n", err)
		http.Error(w, "there was an error when retrieving version", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, version, 200)
}
