package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/resources"
)

// AccessTokenCreate creates a new AccessToken
func AccessTokenCreate(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return
	}

	token, err := resources.CreateRandomAccessToken(db)
	// @todo handle failed save w/out error?
	if err != nil {
		log.Printf("Error when creating access token: '%s'\n", err)
		http.Error(w, "there was an error when creating access token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusCreated)
}
