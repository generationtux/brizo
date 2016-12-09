package auth

import (
	"database/sql/driver"
	"testing"

	"github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestIsFirstUserReturnsTrueWhenNoUsers(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		return testdb.RowsFromCSVString([]string{}, ""), nil
	})

	assert.Equal(t, true, IsFirstUser(db))
}

func TestIsFirstUserReturnsFalseWhenUsersExist(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		columns := []string{"username"}
		return testdb.RowsFromCSVString(columns, "foobar"), nil
	})

	assert.Equal(t, false, IsFirstUser(db))
}

func TestGithubUserIsNotAllowed(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	var sentQuery string

	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		sentQuery = query
		return testdb.RowsFromCSVString([]string{}, ""), nil
	})

	assert.Equal(t, false, GithubUserAllowed(db, "foobar"))
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

	assert.Equal(t, true, GithubUserAllowed(db, "foobar"))
	assert.Equal(t, "SELECT * FROM \"users\"  WHERE (github_username = ?)", sentQuery)
}
