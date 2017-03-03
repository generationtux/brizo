package resources

import (
	"database/sql/driver"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateEnvironmentConfig(t *testing.T) {
	config := EnvironmentConfig{
		Name:            "foo",
		Value:           "bar",
		EnvironmentUUID: "1",
	}
	db, _ := gorm.Open("testdb", "")
	var query string
	var args []driver.Value

	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return driver.ResultNoRows, nil
	})

	CreateEnvironmentConfig(db, &config)
	expectQuery := "INSERT INTO \"environment_configs\" (\"created_at\",\"updated_at\",\"name\",\"value\",\"environment_uuid\") VALUES (?,?,?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, "foo", args[2])
	assert.Equal(t, "bar", args[3])
}
