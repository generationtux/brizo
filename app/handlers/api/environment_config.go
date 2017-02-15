package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/resources"
)

// EnvironmentConfigCreate endpoint for creating config entry
func EnvironmentConfigCreate(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	var createForm struct {
		Name          string
		Value         string
		EnvironmentID string
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&createForm)
	defer r.Body.Close()
	if err != nil {
		log.Printf("decoding error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to parse the form", http.StatusInternalServerError)
		return
	}

	envID, _ := strconv.ParseUint(createForm.EnvironmentID, 10, 0)
	envConfg := resources.EnvironmentConfig{
		Name:          createForm.Name,
		Value:         createForm.Value,
		EnvironmentID: envID,
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
