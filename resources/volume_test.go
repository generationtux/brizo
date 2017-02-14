package resources

import (
	"database/sql/driver"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateAVolume(t *testing.T) {
	id := uuid.New()
	volume := Volume{
		Name:   "foobar",
		UUID:   id,
		Type:   "Secret",
		Source: "over there",
	}
	db, _ := gorm.Open("testdb", "")
	var query string
	var args []driver.Value

	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return driver.ResultNoRows, nil
	})

	CreateVolume(db, &volume)
	expectQuery := "INSERT INTO \"volumes\" (\"created_at\",\"updated_at\",\"uuid\",\"name\",\"type\",\"source\",\"version_id\") VALUES (?,?,?,?,?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, id, args[2])
	assert.Equal(t, "foobar", args[3])
	assert.Equal(t, "Secret", args[4])
	assert.Equal(t, "over there", args[5])
}
