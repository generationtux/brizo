package api

import (
	"log"
	"net/http"

	"github.com/generationtux/brizo/app/handlers/jsonutil"
	"github.com/generationtux/brizo/auth"
	"github.com/generationtux/brizo/database"
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
