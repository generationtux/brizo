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
	id := uuid.New()
	app := Application{Slug: "my-app", UUID: "app-uuid123"}
	environment := Environment{Slug: "dev", UUID: "env-uuid123"}
	version := Version{
		Name:        "foobar",
		Slug:        "foobar",
		UUID:        id,
		Application: app,
		Environment: &environment,
	}

	mockClient := new(mockKubeClient)
	expectedError := errors.New("foo error")
	mockClient.On("CreateOrUpdateDeployment", mock.Anything).Return(expectedError)

	db, _ := gorm.Open("testdb", "")
	result, err := CreateVersion(db, mockClient, &version)
	assert.False(t, result)
	assert.Equal(t, "foo error", err.Error())

	mockClient.AssertExpectations(t)
}

func TestCanCreateAVersion(t *testing.T) {
	id := uuid.New()
	app := Application{Slug: "my-app", UUID: "app-uuid123"}
	environment := Environment{Slug: "dev", UUID: "env-uuid123"}
	version := Version{
		Name:        "foobar",
		Slug:        "foobar",
		UUID:        id,
		Application: app,
		Environment: &environment,
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
	expectQuery := "INSERT INTO \"applications\" (\"created_at\",\"updated_at\",\"uuid\",\"name\",\"slug\") VALUES (?,?,?,?,?)"
	assert.Equal(t, expectQuery, query)
	assert.Equal(t, "app-uuid123", args[2])
	assert.Equal(t, "my-app", args[4])

	deployment := versionDeploymentDefinition(&version)
	expectSpec, err := json.Marshal(deployment)
	assert.Nil(t, err)
	assert.Equal(t, string(expectSpec), version.Spec)
}

func TestVersionDeploymentDefinition(t *testing.T) {
	app := Application{Slug: "my-app", UUID: "app-uuid123"}
	environment := Environment{Slug: "dev", UUID: "env-uuid123", Application: app}
	environment.Version.Application = app
	version := &Version{
		Application: app,
		Environment: &environment,
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

	expectDeploymentLabels := map[string]string{"brizoManaged": "true", "appUUID": "app-uuid123", "envUUID": "env-uuid123", "versionUUID": "version-uuid123"}
	assert.Equal(t, expectDeploymentLabels, deployment.Labels)

	expectSelector := map[string]string{"appUUID": "app-uuid123", "envUUID": "env-uuid123", "versionUUID": "version-uuid123"}
	assert.Equal(t, expectSelector, deployment.Spec.Selector.MatchLabels)

	expectTemplateLabels := map[string]string{"brizoManaged": "true", "appUUID": "app-uuid123", "envUUID": "env-uuid123", "versionUUID": "version-uuid123"}
	assert.Equal(t, expectTemplateLabels, deployment.Spec.Template.Labels)
}

func TestRawArgsToParsedArgs(t *testing.T) {
	replacements := map[string]string{
		"${BRIZO_ENVIRONMENT}": "Current-environment",
	}
	expectedEnv := []string{"--hostname=current-environment.payments", "--example=BRIZO_ENVIRONMENT"}
	validEnv := []string{"--hostname=${BRIZO_ENVIRONMENT}.payments", "--example=BRIZO_ENVIRONMENT"}
	missingBracesEnv := []string{"--hostname=$BRIZO_ENVIRONMENT.payments", "--example=BRIZO_ENVIRONMENT"}
	curvedBracesEnv := []string{"--hostname=$(BRIZO_ENVIRONMENT).payments", "--example=BRIZO_ENVIRONMENT"}
	brizoEnvArgs := [][]string{
		validEnv,
		missingBracesEnv,
		curvedBracesEnv,
	}

	assert.Equal(t, expectedEnv, parseRawArguments(brizoEnvArgs[0], replacements))
	assert.Equal(t, brizoEnvArgs[1], parseRawArguments(brizoEnvArgs[1], replacements))
	assert.Equal(t, brizoEnvArgs[2], parseRawArguments(brizoEnvArgs[2], replacements))
}

func TestRawArgumentsToJSON(t *testing.T) {
	version := &Version{
		Containers: []Container{
			Container{Name: "container-1", Args: []string{"echo", "${BRIZO_ENVIRONMENT}"}},
			Container{Name: "container-2", Args: []string{"echo", "Hello,", "World!"}},
		},
	}
	expected := "{\"containers\":[{\"arguments\":[\"echo\",\"${BRIZO_ENVIRONMENT}\"],\"name\":\"container-1\"},{\"arguments\":[\"echo\",\"Hello,\",\"World!\"],\"name\":\"container-2\"}]}"
	hydrateRawArgumentsToJSON(version)

	assert.JSONEq(t, expected, version.RawArguments)
}

func TestRawArgumentJSONToSingleContainerArgs(t *testing.T) {
	version := &Version{
		RawArguments: "{\"containers\":[{\"arguments\":[\"echo\",\"${BRIZO_ENVIRONMENT}\"],\"name\":\"container-1\"},{\"arguments\":[\"echo\",\"Hello,\",\"World!\"],\"name\":\"container-2\"}]}",
	}

	expected := []string{"echo", "${BRIZO_ENVIRONMENT}"}
	actual := rawArgumentJSONToContainerArgs("container-1", version)
	assert.Equal(t, expected, actual)

	expected = []string{}
	actual = rawArgumentJSONToContainerArgs("container-0", version)
	assert.Equal(t, expected, actual)
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

func (m *mockKubeClient) GetService(namespace string, name string) (*v1.Service, error) {
	return &v1.Service{}, nil
}

func (m *mockKubeClient) CreateService(service *v1.Service) error {
	args := m.Called(service)
	return args.Error(0)
}

func (m *mockKubeClient) UpdateService(service *v1.Service) error {
	return nil
}
