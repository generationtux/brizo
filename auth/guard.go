package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/resources"
)

// APIMiddleware authenticates API requests
func APIMiddleware(rw http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		rw.WriteHeader(401)
		rw.Write([]byte("not authorized"))
		return
	}

	if !validateAuthHeader(authHeader) {
		log.Println("Authorization header is invalid:", authHeader)
		rw.WriteHeader(401)
		rw.Write([]byte("invalid authorization header"))
		return
	}

	jwtToken := extractBearerTokenFromHeader(authHeader)
	if ValidateJWTToken(jwtToken) || ValidatePersonalAccessToken(jwtToken) {
		next(rw, request)
		return
	}

	rw.WriteHeader(401)
	rw.Write([]byte("not authorized"))
}

// ValidatePersonalAccessToken check that a valid access token exists in the
// database.
func ValidatePersonalAccessToken(token string) bool {
	// @todo this db connection stuff's is a real mess. This should be kicked into
	// callback or something to allow for any db connections to be caught?
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		// http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
		return false
	}
	if resources.HasAccessToken(db, token) {
		return true
	}
	return false
}

// validateAuthHeader will ensure the header is formatted correctly
// header should be formatted as:
//	- Authorization: Bearer [JWT_TOKEN]
func validateAuthHeader(header string) bool {
	data := []byte(header)
	matched, _ := regexp.Match("^Bearer .*", data)

	return matched
}

// ValidateJWTToken determines if the provided token is valid
func ValidateJWTToken(token string) bool {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret()), nil
	})

	if err != nil {
		log.Println(err.Error()+":", token)
		return false
	}

	return t.Valid
}

// jwtSecret get the configured JWT secret
func jwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

// extractBearerTokenFromHeader will parse the header string and return just the token
// header should be formatted as Bearer [TOKEN]
func extractBearerTokenFromHeader(header string) string {
	return header[7:]
}
