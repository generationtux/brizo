package api

import (
	"log"
	"net/http"

	"github.com/Machiel/slugify"
	"github.com/generationtux/brizo/app/handlers/jsonutil"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/generationtux/brizo/resources"
	"github.com/go-zoo/bone"
	"github.com/pborman/uuid"
)

// DeployCreate deploys a Version
func DeployCreate(w http.ResponseWriter, r *http.Request) {
	var deployForm struct {
		// id of the Environment where we will deploy the version
		EnvironmentUUID string `json:"environment_id"`
	}
	err := jsonRequest(r, &deployForm)
	if err != nil {
		jre := jsonutil.NewJSONResponseError(
			http.StatusBadRequest,
			err.Error())
		jsonutil.RespondJSONError(w, jre)
		return
	}

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"Database connection error")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	deployToEnvironment, err := resources.GetEnvironment(db, deployForm.EnvironmentUUID)
	if err != nil {
		jre := jsonutil.NewJSONResponseError(
			http.StatusNotFound,
			"Environment not found")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	client, err := kube.New()
	if err != nil {
		log.Printf("Kube client error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusServiceUnavailable,
			"Kube connection error")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	version, err := resources.GetVersion(db, bone.GetValue(r, "version-uuid"), client)
	if err != nil {
		jre := jsonutil.NewJSONResponseError(
			http.StatusNotFound,
			"Version not found")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	newVersion := resources.Version{
		UUID:          uuid.New(),
		Name:          version.Name,
		Slug:          slugify.Slugify(version.Name),
		Containers:    version.Containers,
		Replicas:      version.Replicas,
		Volumes:       version.Volumes,
		EnvironmentID: deployToEnvironment.ID,
		Environment:   *deployToEnvironment,
	}

	_, err = resources.CreateVersion(db, client, &newVersion)
	if err != nil {
		log.Printf("Error deploying version: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"Unable to deploy version")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	jsonResponse(w, newVersion, 201)
}
