package api

import (
	"log"
	"net/http"

	"github.com/generationtux/brizo/auth"
	"github.com/generationtux/brizo/database"
	githuboauth "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const htmlIndex = `<html><body>
Log in with <a href="/o/auth/login/github">GitHub</a>
</body></html>
`

var (
	oauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user", "user:email", "repo"},
		Endpoint:     github.Endpoint,
	}
	oauthStateString = auth.GetOAuthStateString()
)

// AuthMainHandler for basic login page
func AuthMainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

// AuthGithubHandler for requesting oauth access from Github
func AuthGithubHandler(w http.ResponseWriter, r *http.Request) {
	auth.HydrateOAuthConfig(oauthConf)
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// AuthGithubCallbackHandler for handling oauth access response from Github
func AuthGithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	auth.HydrateOAuthConfig(oauthConf)
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := githuboauth.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Printf("client.Users.Get() faled with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	auth.CreateNewGithubUser(db, user, token.AccessToken)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
