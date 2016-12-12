package api

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

const htmlIndex = `<!doctype html><html><body>
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
		http.Error(w, "there was an error when attempting to connect to the database", http.StatusInternalServerError)
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
		http.Error(w, "there was an error when attemping to create a new user", http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// AuthGithubHandler for requesting oauth access from Github
func AuthGithubHandler(w http.ResponseWriter, r *http.Request) {
	auth.HydrateOAuthConfig(oauthConf)
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// AuthGithubHandler for handling oauth access response from Github
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

	if auth.IsFirstUser(db) {
		auth.CreateNewGithubUser(db, user, token.AccessToken)
		http.Redirect(w, r, "/get-token", http.StatusTemporaryRedirect)
		return
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
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// @todo user is not allowed
	log.Printf("'%s' attempted to authenticate without an account\n", *user.Login)
	http.Error(w, "user has not been created prior to authentication", http.StatusForbidden)
}
