package resources

import (
	"database/sql/driver"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateAVersion(t *testing.T) {
	id := uuid.New()
	version := Version{
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

	CreateVersion(db, &version)
	expectQuery := "INSERT INTO \"versions\" (\"created_at\",\"updated_at\",\"uuid\",\"name\",\"slug\",\"image\",\"environment_id\") VALUES (?,?,?,?,?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, id, args[2])
	assert.Equal(t, "foobar", args[3])
}

func TestVersionDeploymentDefinition(t *testing.T) {
	app := Application{Slug: "my-app", UUID: "app-uuid123"}
	environment := Environment{Slug: "dev", UUID: "env-uuid123"}
	environment.Application = app
	version := &Version{
		Environment: environment,
		Name:        "version 1",
		Slug:        "version-1",
		Image:       "foo:latest",
		UUID:        "version-uuid123",
		Replicas:    3,
	}

	deployment := versionDeploymentDefinition(version)
	assert.Equal(t, "my-app-dev-version-1", deployment.Name)
	assert.Equal(t, "brizo", deployment.Namespace)
	assert.Equal(t, int32(version.Replicas), *deployment.Spec.Replicas)
	assert.Equal(t, "foo:latest", deployment.Spec.Template.Spec.Containers[0].Image)

	expectDeploymentLabels := map[string]string{"brizoManaged": "true", "appUUID": "app-uuid123", "envUUID": "env-uuid123"}
	assert.Equal(t, expectDeploymentLabels, deployment.Labels)

	expectSelector := map[string]string{"envUUID": "env-uuid123", "versionUUID": "version-uuid123"}
	assert.Equal(t, expectSelector, deployment.Spec.Selector.MatchLabels)

	expectTemplateLabels := map[string]string{"brizoManaged": "true", "appUUID": "app-uuid123", "envUUID": "env-uuid123", "versionUUID": "version-uuid123"}
	assert.Equal(t, expectTemplateLabels, deployment.Spec.Template.Labels)
}
