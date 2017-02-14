package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/generationtux/brizo/app/handlers/jsonutil"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/resources"
	"github.com/go-zoo/bone"
)

// VolumeIndex provides a listing of all Volumes
func VolumeIndex(w http.ResponseWriter, r *http.Request) {
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

	volumes, err := resources.GetVolumesByVersionUUID(db, bone.GetValue(r, "version-uuid"))
	if err != nil {
		log.Printf("Error when retrieving volumes: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to connect to the database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(volumes)
}

// VolumeCreate creates a new Volume
func VolumeCreate(w http.ResponseWriter, r *http.Request) {
	// client, err := kube.New()
	// if err != nil {
	// 	log.Printf("Kube client error: '%s'\n", err)
	// 	http.Error(w, "unable to reach Kubernetes", http.StatusInternalServerError)
	// 	return
	// }

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	var createForm struct {
		Name   string
		Type   string
		Source string
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&createForm)
	defer r.Body.Close()
	if err != nil {
		log.Printf("decoding error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to parse the form", http.StatusInternalServerError)
		return
	}

	version, err := resources.GetVersion(db, bone.GetValue(r, "version-uuid"))
	if err != nil {
		log.Printf("Error when retrieving version: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving version")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	volume := resources.Volume{
		Name:      createForm.Name,
		Type:      createForm.Type,
		Source:    createForm.Source,
		VersionID: version.ID,
		// Version:   *version,
	}
	_, err = resources.CreateVolume(db, &volume)
	// @todo handle failed save w/out error?
	if err != nil {
		log.Printf("Error when retrieving volume: '%s'\n", err)
		http.Error(w, "there was an error when retrieving volume", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, volume, 201)
}
