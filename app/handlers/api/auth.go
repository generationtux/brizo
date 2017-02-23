package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/generationtux/brizo/app/handlers/jsonutil"
	"github.com/generationtux/brizo/auth"
	"github.com/generationtux/brizo/database"
	"github.com/go-zoo/bone"
	"github.com/mholt/binding"
)

// AuthCreateUser will invite a new user by Github username
func AuthCreateUser(w http.ResponseWriter, r *http.Request) {
	createUserForm := new(auth.CreateUserForm)
	errs := binding.Bind(r, createUserForm)
	if errs.Handle(w) {
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
	user := auth.User{
		Username:       createUserForm.Username,
		GithubUsername: createUserForm.Username,
	}
	success, err := auth.CreateUser(db, &user)
	if success == false {
		log.Printf("%s when trying to create user\n", err)
		// @todo correctly redirect w/ bad request shown
		jre := jsonutil.NewJSONResponseError(
			http.StatusBadRequest,
			"there was an error when attemping to create a new user")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return
}

// ValidateToken validates the provided JWT token
func ValidateToken(rw http.ResponseWriter, request *http.Request) {
	tokenRequest := new(auth.ValidateJWTForm)
	errs := binding.Bind(request, tokenRequest)
	if errs.Handle(rw) {
		return
	}

	if !auth.ValidateJWTToken(tokenRequest.Token) {
		rw.Write([]byte("invalid token"))
		rw.WriteHeader(401)
	}

	rw.Write([]byte("ok"))
	rw.WriteHeader(200)
}

// AuthGetInvitees provides a listing of users who have not yet accepted their
// invitation.
func AuthGetInvitees(w http.ResponseWriter, r *http.Request) {
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
	invitees, err := auth.GetInvitedUsers(db)
	if err != nil {
		log.Printf("Error when retrieving invitees: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to connect to the database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(invitees)
}

// AuthDeleteInvitees deletes an unaccepted invite
func AuthDeleteInvitees(w http.ResponseWriter, r *http.Request) {
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
	id := bone.GetValue(r, "id")
	success, err := auth.DeleteInvitedUser(db, id)
	if err != nil {
		log.Printf("Error when deleting invitee: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to connect to the database")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	if !success {
		log.Printf("Unsuccesful when attempting to delete invitee: '%s'\n", err)
		jre := jsonutil.NewJSONResponseError(
			http.StatusInternalServerError,
			"there was an error when attempting to delete the invitee")
		jsonutil.RespondJSONError(w, jre)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
