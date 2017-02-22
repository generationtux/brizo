package resources

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func TestCreateVersionDoesNotStoreWhenKubeFails(t *testing.T) {
	version := &Version{Name: "foo"}
	mockClient := new(mockKubeClient)
	expectedError := errors.New("foo error")
	mockClient.On("CreateOrUpdateDeployment", mock.Anything).Return(expectedError)

	db, _ := gorm.Open("testdb", "")
	result, err := CreateVersion(db, mockClient, version)
	assert.False(t, result)
	assert.Equal(t, "foo error", err.Error())

	mockClient.AssertExpectations(t)
}

func TestCanCreateAVersion(t *testing.T) {
	id := uuid.New()
	version := Version{
		Name: "foobar",
		Slug: "foobar",
		UUID: id,
	}

	mockClient := new(mockKubeClient)
	mockClient.On("CreateOrUpdateDeployment", mock.Anything).Return(nil)

	db, _ := gorm.Open("testdb", "")
	var query string
	var args []driver.Value

	testdb.SetExecWithArgsFunc(func(q string, a []driver.Value) (driver.Result, error) {
		query = q
		args = a
		return driver.ResultNoRows, nil
	})

	CreateVersion(db, mockClient, &version)
	expectQuery := "INSERT INTO \"versions\" (\"created_at\",\"updated_at\",\"uuid\",\"name\",\"slug\",\"environment_id\",\"spec\") VALUES (?,?,?,?,?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, id, args[2])
	assert.Equal(t, "foobar", args[3])

	deployment := versionDeploymentDefinition(&version)
	expectSpec, err := json.Marshal(deployment)
	assert.Nil(t, err)
	assert.Equal(t, string(expectSpec), args[6])
}

func TestVersionDeploymentDefinition(t *testing.T) {
	app := Application{Slug: "my-app", UUID: "app-uuid123"}
	environment := Environment{Slug: "dev", UUID: "env-uuid123"}
	environment.Application = app
	version := &Version{
		Environment: environment,
		Name:        "version 1",
		Slug:        "version-1",
		UUID:        "version-uuid123",
		Replicas:    3,
		Containers: []Container{
			Container{
				Name:  "app",
				Image: "foo:latest",
			},
		},
	}

	deployment := versionDeploymentDefinition(version)
	assert.Equal(t, "my-app-dev", deployment.Name)
	assert.Equal(t, "brizo", deployment.Namespace)
	assert.Equal(t, int32(version.Replicas), *deployment.Spec.Replicas)
	assert.Equal(t, "app", deployment.Spec.Template.Spec.Containers[0].Name)
	assert.Equal(t, "foo:latest", deployment.Spec.Template.Spec.Containers[0].Image)

	expectDeploymentLabels := map[string]string{"brizoManaged": "true", "appUUID": "app-uuid123", "envUUID": "env-uuid123"}
	assert.Equal(t, expectDeploymentLabels, deployment.Labels)

	expectSelector := map[string]string{"envUUID": "env-uuid123", "versionUUID": "version-uuid123"}
	assert.Equal(t, expectSelector, deployment.Spec.Selector.MatchLabels)

	expectTemplateLabels := map[string]string{"brizoManaged": "true", "appUUID": "app-uuid123", "envUUID": "env-uuid123", "versionUUID": "version-uuid123"}
	assert.Equal(t, expectTemplateLabels, deployment.Spec.Template.Labels)
}

// Mock kube client

type mockKubeClient struct {
	mock.Mock
}

func (m *mockKubeClient) Health() error {
	return nil
}

func (m *mockKubeClient) GetPods(name string, options v1.ListOptions) ([]v1.Pod, error) {
	return []v1.Pod{}, nil
}

func (m *mockKubeClient) CreateDeployment(deployment *v1beta1.Deployment) error {
	args := m.Called(deployment)
	return args.Error(0)
}

func (m *mockKubeClient) DeleteDeployment(deployment *v1beta1.Deployment) error {
	return nil
}

func (m *mockKubeClient) CreateOrUpdateDeployment(deployment *v1beta1.Deployment) error {
	args := m.Called(deployment)
	return args.Error(0)
}

func (m *mockKubeClient) FindDeploymentByName(namespace, name string) (*v1beta1.Deployment, error) {
	return &v1beta1.Deployment{}, nil
}

func (m *mockKubeClient) GetServices(name string, options v1.ListOptions) ([]v1.Service, error) {
	return []v1.Service{}, nil
}

func (m *mockKubeClient) CreateService(service *v1.Service) error {
	args := m.Called(service)
	return args.Error(0)
}
