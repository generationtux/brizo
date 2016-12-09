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
	user := auth.User{
		Username:       createUserForm.Username,
		GithubUsername: createUserForm.Username,
	}
	success, err := auth.CreateUser(&user)
	if success == false {
		log.Printf("%s when trying to create user\n", err)
		// @todo correctly redirect w/ bad request shown
		http.Redirect(w, r, "/", http.StatusBadRequest)
	} else {
		// @todo correctly redirect w/ created shown
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token, err := oauthConf.Exchange(oauth2.NoContext, oAuthCallbackForm.Code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := githuboauth.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Printf("client.Users.Get() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	db, _ := database.Connect()
	defer db.Close()
	if auth.IsFirstUser(db) {
		auth.CreateNewGithubUser(user, token.AccessToken)
	} else if auth.GithubUserAllowed(db, *user.Login) {
		brizoUser := auth.User{
			Username:       *user.Login,
			Name:           *user.Name,
			Email:          *user.Email,
			GithubUsername: *user.Login,
			GithubToken:    token.AccessToken,
		}
		success, err := auth.UpdateUser(&brizoUser)
		if success == false {
			log.Printf("failed to update user '%s' because '%s'\n", *user.Login, err)
			// @todo update to correct status code
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
		// @todo maybe add the correct status code here too
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
