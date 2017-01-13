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

var (
	oauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user", "user:email", "repo"},
		Endpoint:     github.Endpoint,
	}
	oauthStateString = auth.GetOAuthStateString()
)

// AuthGithubHandler for requesting oauth access from Github
func AuthGithubHandler(w http.ResponseWriter, r *http.Request) {
	auth.HydrateOAuthConfig(oauthConf)
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func authErrorRedirect(w http.ResponseWriter, r *http.Request) {
	endpoint := "/app/login?err=1"
	http.Redirect(w, r, endpoint, http.StatusTemporaryRedirect)
}
func authDenyRedirect(w http.ResponseWriter, r *http.Request) {
	endpoint := "/app/login?err=2"
	http.Redirect(w, r, endpoint, http.StatusTemporaryRedirect)
}

// AuthGithubCallbackHandler for handling oauth access response from Github
// any errors here should return a redirect back to the apps login form to display the error
func AuthGithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	auth.HydrateOAuthConfig(oauthConf)
	oAuthCallbackForm := new(auth.OAuthCallbackForm)
	errs := binding.Bind(r, oAuthCallbackForm)
	if errs.Len() > 0 {
		authErrorRedirect(w, r)
		return
	}

	if oAuthCallbackForm.State != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, oAuthCallbackForm.State)
		authErrorRedirect(w, r)
		return
	}

	token, err := oauthConf.Exchange(oauth2.NoContext, oAuthCallbackForm.Code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		authErrorRedirect(w, r)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := githuboauth.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Println("unable to get user details")
		authErrorRedirect(w, r)
		return
	}

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Println(err)
		authErrorRedirect(w, r)
		return
	}

	var jwtToken string
	var jwtError error
	if auth.IsFirstUser(db) {
		user, err := auth.CreateNewGithubUser(db, user, token.AccessToken)

		if err != nil {
			log.Println(err)
			authErrorRedirect(w, r)
			return
		}

		jwtToken, jwtError = auth.CreateJWTToken(user)
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

		if success == false || err != nil {
			authErrorRedirect(w, r)
			return
		}

		jwtToken, jwtError = auth.CreateJWTToken(brizoUser)
	} else {
		// user is not allowed
		authDenyRedirect(w, r)
	}

	if jwtError != nil {
		log.Println(jwtError)
		authErrorRedirect(w, r)
	}

	http.Redirect(w, r, "/app/auth?token="+jwtToken, http.StatusTemporaryRedirect)
	return
}
