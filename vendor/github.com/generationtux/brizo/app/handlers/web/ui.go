package web

import (
	"log"
	"net/http"
	"time"
    "fmt"

	"github.com/CloudyKit/jet"
	jwt "github.com/dgrijalva/jwt-go"
    _ "github.com/auth0/go-jwt-middleware"
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
	expireToken := time.Now().Add(time.Hour * 3).Unix()

	// create claims
	claims := &jwt.StandardClaims{
		ExpiresAt: expireToken,
		Issuer:    "brizo",
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := token.SignedString(mySigningKey)

    rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Authorization", fmt.Sprintf("Bearer %v", signedToken))

	// Redirect the user to his profile
	//http.Redirect(rw, request, "/app", http.StatusTemporaryRedirect)
}

/*
var jwtMiddlewareHandler = jwtMiddleware.New(jwtMiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
        decoded, err := base64.URLEncoding.DecodeString(os.Getenv("AUTH0_CLIENT_SECRET"))
        if err != nil {
            return nil, err
        }
        return decoded, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})
*/
