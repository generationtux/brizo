package auth

import (
	"testing"

	"database/sql/driver"
	"github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

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
