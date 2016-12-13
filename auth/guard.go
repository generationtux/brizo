package auth

import (
	"log"
	"net/http"
	"regexp"
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

	if !validateJWTToken(authHeader) {
		rw.WriteHeader(401)
		rw.Write([]byte("not authorized"))
		return
	}

	next(rw, request)
}

// validateAuthHeader will ensure the header is formatted correctly
// header should be formatted as:
//	- Authorization: Bearer [JWT_TOKEN]
func validateAuthHeader(header string) bool {
	data := []byte(header)
	matched, _ := regexp.Match("^Bearer .*", data)

	return matched
}

// validateJWTToken @todo in progress from Lark
func validateJWTToken(token string) bool {
	return true
}
