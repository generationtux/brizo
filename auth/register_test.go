package auth

import (
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/erikstmartin/go-testdb"
	"github.com/generationtux/brizo/config"
	githuboauth "github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestAuthStateStringGeneration(t *testing.T) {
	stateString := GetOAuthStateString()
	matches, _ := regexp.Match("[a-zA-Z0-9]", []byte(stateString))

	assert.Len(t, stateString, 64)
	assert.True(t, true, matches)
}

func TestCreateNewGithubUser(t *testing.T) {
	token := "123"
	login := "fooUser"
	name := "Foo Name"
	email := "foo@email.com"
	user := &githuboauth.User{
		Login: &login,
		Name:  &name,
		Email: &email,
	}

	db, _ := gorm.Open("testdb", "")
	var query string
	var args []driver.Value

	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return driver.ResultNoRows, nil
	})

	CreateNewGithubUser(db, user, token)
	expectQuery := "INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"username\",\"name\",\"email\",\"github_username\",\"github_token\") VALUES (?,?,?,?,?,?,?)"
	assert.Equal(t, expectQuery, query)
	// args are prefixed with created_at and update_at
	assert.Equal(t, login, args[2])
	assert.Equal(t, name, args[3])
	assert.Equal(t, email, args[4])
	assert.Equal(t, login, args[5])
}

func TestHydratesConfig(t *testing.T) {
	config.App.OAuthGithubClientID = "exampleId"
	config.App.OAuthGithubClientSecret = "exampleSecret"
	oauthConf := oauth2.Config{}

	assert.Equal(t, "", oauthConf.ClientID)
	assert.Equal(t, "", oauthConf.ClientSecret)
	HydrateOAuthConfig(&oauthConf)
	assert.Equal(t, "exampleId", oauthConf.ClientID)
	assert.Equal(t, "exampleSecret", oauthConf.ClientSecret)
}

func TestIsFirstUserReturnsTrueWhenNoUsers(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		return testdb.RowsFromCSVString([]string{}, ""), nil
	})

	assert.Equal(t, true, isFirstUser(db))
}

func TestIsFirstUserReturnsFalseWhenUsersExist(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		columns := []string{"username"}
		return testdb.RowsFromCSVString(columns, "foobar"), nil
	})

	assert.Equal(t, false, isFirstUser(db))
}

func TestGithubUserIsNotAllowed(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	var sentQuery string

	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		sentQuery = query
		return testdb.RowsFromCSVString([]string{}, ""), nil
	})

	assert.Equal(t, false, githubUserAllowed(db, "foobar"))
	assert.Equal(t, "SELECT * FROM \"users\"  WHERE (github_username = ?)", sentQuery)
}

func TestGithubUserIsAllowed(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	var sentQuery string

	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		sentQuery = query
		columns := []string{"github_username"}
		return testdb.RowsFromCSVString(columns, "foobar"), nil
	})

	assert.Equal(t, true, githubUserAllowed(db, "foobar"))
	assert.Equal(t, "SELECT * FROM \"users\"  WHERE (github_username = ?)", sentQuery)
}
