package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/dgrijalva/jwt-go"
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

	if !ValidateJWTToken(authHeader) {
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

// ValidateJWTToken determines if the provided token is valid
func ValidateJWTToken(token string) bool {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return getJWTSecret(), nil
	})

	return t.Valid
}

func getJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

func getJWTSigningMethod() string {
	return os.Getenv("JWT_SIGNING")
}
