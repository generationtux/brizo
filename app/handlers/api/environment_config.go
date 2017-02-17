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

// GetEnvironmentConfig endpoint for getting configuration
func GetEnvironmentConfig(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(http.StatusInternalServerError, "unable to connect to database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	id := bone.GetValue(r, "environment-uuid")
	config, err := resources.GetEnvironmentConfig(db, id)
	if err != nil {
		log.Printf("Error when retrieving configuration: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when retrieving configuration")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// CreateEnvironmentConfig endpoint for creating config entry
func CreateEnvironmentConfig(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	var createForm struct {
		Name            string
		Value           string
		EnvironmentUUID string
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&createForm)
	defer r.Body.Close()
	if err != nil {
		log.Printf("decoding error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to parse the form", http.StatusInternalServerError)
		return
	}

	envConfg := resources.EnvironmentConfig{
		Name:            createForm.Name,
		Value:           createForm.Value,
		EnvironmentUUID: createForm.EnvironmentUUID,
	}

	_, err = resources.CreateEnvironmentConfig(db, &envConfg)
	// @todo handle failed save w/out error?
	if err != nil {
		log.Printf("Error when retrieving environment config: '%s'\n", err)
		http.Error(w, "there was an error when retrieving environment config", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, envConfg, 200)
}

// DeleteEnvironmentConfig endpoint for deleting config entry
func DeleteEnvironmentConfig(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	var deleteForm struct {
		configID string
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&deleteForm)
	defer r.Body.Close()
	if err != nil {
		log.Printf("decoding error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to parse the form", http.StatusInternalServerError)
		return
	}

	configID, _ := strconv.ParseUint(deleteForm.configID, 10, 0)

	resources.DeleteEnvironmentConfig(db, configID)

	jsonResponse(w, deleteForm, 204)
}
