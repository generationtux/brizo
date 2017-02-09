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

	emails, _, err := client.Users.ListEmails(&githuboauth.ListOptions{})
	var email string
	for _, userEmail := range emails {
		if *userEmail.Primary {
			email = *userEmail.Email
		}
	}

	if err != nil {
		log.Println("unable to get user email")
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
		brizoUser, err := auth.CreateNewGithubUser(db, user, email, token.AccessToken)

		if err != nil {
			log.Println(err)
			authErrorRedirect(w, r)
			return
		}

		jwtToken, jwtError = auth.CreateJWTToken(brizoUser)
	} else if auth.GithubUserAllowed(db, *user.Login) {
		brizoUser := auth.BuildUserFromGithubUser(user, email, token.AccessToken)
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
