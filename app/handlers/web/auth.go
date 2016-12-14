package web

import (
	"log"
	"net/http"

	"github.com/generationtux/brizo/auth"
	"github.com/generationtux/brizo/database"
	githuboauth "github.com/google/go-github/github"
	"github.com/mholt/binding"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const htmlAddUser = `<!doctype html><html><body>
<form action="/api/v1/users" method="post">
  <label>Add user:</label>
  <input type="text" name="username">
</form>
</body></html>
`

// AuthAddNewUser @todo move to JS UI
func AuthAddNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlAddUser))
}

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

// AuthMainHandler @todo move to JS UI
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
	oAuthCallbackForm := new(auth.OAuthCallbackForm)
	errs := binding.Bind(r, oAuthCallbackForm)
	if errs.Handle(w) {
		return
	}
	if oAuthCallbackForm.State != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, oAuthCallbackForm.State)
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	token, err := oauthConf.Exchange(oauth2.NoContext, oAuthCallbackForm.Code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Error(w, "there was an error when establishing an oauth exchange", http.StatusBadRequest)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := githuboauth.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Printf("client.Users.Get() failed with '%s'\n", err)
		http.Error(w, "there was an error while attempting to get user details", http.StatusBadGateway)
		return
	}
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
		http.Error(w, "there was an error while connecting to the database", http.StatusInternalServerError)
		return
	}

	var jwtToken string
	var jwtError error
	if auth.IsFirstUser(db) {
		user, err := auth.CreateNewGithubUser(db, user, token.AccessToken)

		if err != nil {
			log.Printf("failed to create user '%s' because '%s'\n", *user.Login, err)
			http.Error(w, "there was an error when creating new user", http.StatusInternalServerError)
			return
		}

		jwtToken, jwtError = auth.CreateToken(user)
	} else if auth.GithubUserAllowed(db, *user.Login) {
		// @todo check that non-required attributes exist
		brizoUser := auth.User{
			Username:       *user.Login,
			Name:           *user.Name,
			Email:          *user.Email,
			GithubUsername: *user.Login,
			GithubToken:    token.AccessToken,
		}

		success, err := auth.UpdateUser(db, &brizoUser)

		if success == false {
			log.Printf("failed to update user '%s' because '%s'\n", *user.Login, err)
			// @todo update to correct status code
			http.Error(w, "there was an error when updating the user", http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Redirect(w, r, "/app/login", http.StatusTemporaryRedirect)
			return
		}

		jwtToken, jwtError = auth.CreateToken(brizoUser)
	}

	if jwtError != nil {
		log.Printf("failed to create jwt token because '%s'\n", jwtError)

		// @todo update to correct status code
		http.Error(w, "there was an error creating jwt token", http.StatusInternalServerError)
		return
	} else {
		http.Redirect(w, r, "/app/auth?token="+jwtToken, http.StatusTemporaryRedirect)
		return
	}

	// @todo user is not allowed
	log.Printf("'%s' attempted to authenticate without an account\n", *user.Login)
	http.Error(w, "user has not been created prior to authentication", http.StatusForbidden)
}
