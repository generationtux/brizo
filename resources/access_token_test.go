package resources

import (
	"database/sql/driver"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateAnAccessToken(t *testing.T) {
	var query string
	var args []driver.Value
	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return driver.ResultNoRows, nil
	})

	token := generateRandomString(PersonalAccessTokenLength)
	accessToken := AccessToken{
		Token: token,
	}
	db, _ := gorm.Open("testdb", "")
	defer db.Close()

	CreateAccessToken(db, &accessToken)

	expectQuery := "INSERT INTO \"access_tokens\" (\"created_at\",\"updated_at\",\"token\") VALUES (?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, token, args[2])
}

func TestCanCheckExistanceOfAnAccessToken(t *testing.T) {
	testdb.SetQueryWithArgsFunc(func(q string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "created_at", "updated_at", "token"}
		var rows string
		if args[0] == "existingToken" {
			rows = "1,2017-02-01 09:00:00,2017-02-01 09:00:00,existingToken"
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	db, _ := gorm.Open("testdb", "")
	defer db.Close()
	found := HasAccessToken(db, "existingToken")
	missing := HasAccessToken(db, "missingToken")

	assert.True(t, found)
	assert.False(t, missing)
}

func TestCanCreateAnRandomAccessToken(t *testing.T) {
	var query string
	var args []driver.Value
	result := testdb.NewResult(1, nil, 1, nil)
	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return result, nil
	})

	db, _ := gorm.Open("testdb", "")
	defer db.Close()
	token, _ := CreateRandomAccessToken(db)

	expectQuery := "INSERT INTO \"access_tokens\" (\"created_at\",\"updated_at\",\"token\") VALUES (?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Len(t, token.Token, PersonalAccessTokenLength)
}
