package web

import (
	"log"
	"net/http"
	"time"

	"github.com/CloudyKit/jet"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/generationtux/brizo/database"
)

var views = jet.NewHTMLSet("./ui")

// this will be moved to .env
var mySigningKey = []byte("secret")

// RootHandler redirects to Javascript app
func RootHandler(rw http.ResponseWriter, request *http.Request) {
	http.Redirect(rw, request, "/app", 301)
}

// UiHandler for requests to Javascript app
func UiHandler(rw http.ResponseWriter, request *http.Request) {
	view, err := views.GetTemplate("index.html")

	if err != nil {
		log.Println("Unexpected template err:", err.Error())
	}

	view.Execute(rw, nil, nil)
}

// HealthzHandler for health check requests
func HealthzHandler(rw http.ResponseWriter, request *http.Request) {
	err := database.Health()

	if err != nil {
		log.Println(err)
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
		return
	}

	rw.Write([]byte("ok"))
}

// JWT Handler for authentication middleware
func GetTokenHandler(rw http.ResponseWriter, request *http.Request) {
	// Expires the token and cookie in 1 hour
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	expireCookie := time.Now().Add(time.Hour * 1)

	// create claims
	claims := &jwt.StandardClaims{
		ExpiresAt: expireToken,
		Issuer:    "brizo",
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := token.SignedString(mySigningKey)

	cookie := http.Cookie{
		Name:     "id_token",
		Value:    signedToken,
		Expires:  expireCookie,
		HttpOnly: true,
	}

	http.SetCookie(rw, &cookie)

	// Redirect the user to his profile
	http.Redirect(rw, request, "/app", http.StatusTemporaryRedirect)
}
