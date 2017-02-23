package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/generationtux/brizo/app/handlers/jsonutil"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/resources"
)

// AccessTokenCreate creates a new AccessToken
func AccessTokenCreate(w http.ResponseWriter, r *http.Request) {
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

	token, err := resources.CreateRandomAccessToken(db)
	// @todo handle failed save w/out error?
	if err != nil {
		log.Printf("Error when creating access token: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when creating access token")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusCreated)
}
