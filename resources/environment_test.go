package resources

import (
	"database/sql/driver"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateAnEnvironment(t *testing.T) {
	id := uuid.New()
	app := Environment{
		Name: "foobar",
		UUID: id,
	}
	db, _ := gorm.Open("testdb", "")
	var query string
	var args []driver.Value

	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return driver.ResultNoRows, nil
	})

	CreateEnvironment(db, &app)
	expectQuery := "INSERT INTO \"environments\" (\"created_at\",\"updated_at\",\"uuid\",\"name\",\"application_id\") VALUES (?,?,?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, id, args[2])
	assert.Equal(t, "foobar", args[3])
}
