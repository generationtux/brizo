package resources

import (
	"database/sql/driver"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateAnApplication(t *testing.T) {
	app := Application{
		Name: "foobar",
	}
	db, _ := gorm.Open("testdb", "")
	var query string
	var args []driver.Value

	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return driver.ResultNoRows, nil
	})

	CreateApplication(db, &app)
	expectQuery := "INSERT INTO \"applications\" (\"created_at\",\"updated_at\",\"name\") VALUES (?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, "foobar", args[2])
}
