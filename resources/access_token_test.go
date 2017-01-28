package resources

import (
	"database/sql/driver"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateAnAccessToken(t *testing.T) {
	token := generateRandomString(PersonalAccessTokenLength)
	accessToken := AccessToken{
		Token: token,
	}
	db, _ := gorm.Open("testdb", "")
	var query string
	var args []driver.Value

	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return driver.ResultNoRows, nil
	})

	CreateAccessToken(db, &accessToken)
	expectQuery := "INSERT INTO \"access_tokens\" (\"created_at\",\"updated_at\",\"token\") VALUES (?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, token, args[2])
}

// func TestCanCheckExistanceOfAnAccessToken(t *testing.T) {
// 	token := generateRandomString(PersonalAccessTokenLength)
// 	db, _ := gorm.Open("testdb", "")
// 	var query string
// 	var args []driver.Value
//
// 	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
// 		query = q
// 		args = a
// 		return driver.ResultNoRows, nil
// 	})
//
// 	HasAccessToken(db, token)
// 	expectQuery := "SELECT * FROM \"access_tokens\" WHERE token = " + token
// 	assert.Equal(t, expectQuery, query)
// }

// func TestCanCreateAnRandomAccessToken(t *testing.T) {
// 	db, _ := gorm.Open("testdb", "")
// 	var query string
// 	var args []driver.Value
//
// 	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
// 		query = q
// 		args = a
// 		return driver.ResultNoRows, nil
// 	})
//
// 	token, _ := CreateRandomAccessToken(db)
// 	expectQuery := "INSERT INTO \"access_tokens\" (\"created_at\",\"updated_at\",\"token\") VALUES (?,?,?)"
// 	assert.Equal(t, expectQuery, query)
// 	assert.Len(t, token.Token, auth.PersonalAccessTokenLength)
// }
