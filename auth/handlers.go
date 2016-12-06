package auth

import (
	"fmt"
	"net/http"

	"github.com/generationtux/brizo/database"
	"github.com/google/go-github/github"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user", "user:email", "repo"},
		Endpoint:     githuboauth.Endpoint,
	}
	oauthStateString = generateRandomString(64)
)

const htmlIndex = `<html><body>
Log in with <a href="/o/auth/login/github">GitHub</a>
</body></html>
`

func AuthMainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

func AuthGithubHandler(w http.ResponseWriter, r *http.Request) {
	hydrateOAuthConfig(oauthConf)
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func AuthGithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	hydrateOAuthConfig(oauthConf)
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		fmt.Printf("client.Users.Get() faled with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
	// fmt.Printf("access token set to '%s'\n", token.AccessToken)
	createNewGithubUser(user)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func createNewGithubUser(githubUser *github.User) {
	user := User{
		Username: *githubUser.Login,
		Name:     *githubUser.Name,
		Email:    *githubUser.Email,
	}
	db, e := database.Connect()
	defer db.Close()
	fmt.Printf("%v\n%v\n", db, e)
	db.Create(&user)
}
