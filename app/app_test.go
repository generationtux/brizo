package app

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApplication(t *testing.T) {
	app := New()
	assert.NotEmpty(t, app.serverListener)
	assert.NotEmpty(t, app.serverHandler)
	assert.NotEmpty(t, app.healthChecks)
	assert.NotEmpty(t, app.migrator)
}

func TestAppServer(t *testing.T) {
	mockServer := func(address string, handler http.Handler) error {
		_, ok := handler.(mockHandler)
		if !ok {
			t.Fatal("Expected handler to be passed to server.")
		}
		return errors.New("custom-server-error")
	}

	app := &Application{serverListener: mockServer, serverHandler: mockHandler{}}
	err := app.Server()
	assert.Equal(t, "custom-server-error", err.Error())
}

type mockHandler struct{}

func (m mockHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func TestAppInitMigrations(t *testing.T) {
	mockMigrator := func() error {
		return errors.New("custom-migrator-error")
	}

	app := &Application{
		healthChecks: []ChecksHealth{},
		migrator:     mockMigrator,
	}

	err := app.Initialize()
	assert.Equal(t, "custom-migrator-error", err.Error())
}

func TestAppInitHealth(t *testing.T) {
	mockMigrator := func() error { return nil }
	app := &Application{healthChecks: []ChecksHealth{}, migrator: mockMigrator}
	assert.Nil(t, app.Initialize())

	checkOne := func() error { return nil }
	checkTwo := func() error { return errors.New("custom-health-error") }
	app = &Application{healthChecks: []ChecksHealth{checkOne, checkTwo}}

	err := app.Initialize()
	assert.Equal(t, "custom-health-error", err.Error())
}

func TestGetAddress(t *testing.T) {
	os.Setenv("BRIZO_PORT", "")
	assert.Equal(t, ":8080", getAddress())

	os.Setenv("BRIZO_PORT", "1234")
	assert.Equal(t, ":1234", getAddress())
}
