package resources

import (
	"database/sql/driver"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCanCreateAnEnvironment(t *testing.T) {
	id := uuid.New()
	env := Environment{
		Name: "foobar",
		UUID: id,
	}

	mockClient := new(mockKubeClient)
	mockClient.On("CreateService", mock.Anything).Return(nil)

	db, _ := gorm.Open("testdb", "")
	var query string
	var args []driver.Value

	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return driver.ResultNoRows, nil
	})

	CreateEnvironment(db, mockClient, &env, &Application{
		Slug: "foo",
		UUID: "abc-abc-abc-abc",
	})
	expectQuery := "INSERT INTO \"environments\" (\"created_at\",\"updated_at\",\"uuid\",\"name\",\"slug\",\"version_id\",\"version_uuid\",\"application_id\",\"application_uuid\") VALUES (?,?,?,?,?,?,?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, id, args[2])
	assert.Equal(t, "foobar", args[3])
}
