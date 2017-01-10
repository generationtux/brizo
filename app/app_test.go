package app

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApplication(t *testing.T) {
	app := New()
	assert.NotEmpty(t, app.version)
	assert.NotEmpty(t, app.server)
	assert.NotEmpty(t, app.runner)
}

func TestAppVersion(t *testing.T) {
	app := &Application{version: "1.x"}
	assert.Equal(t, "1.x", app.Version())
}

func TestAppServer(t *testing.T) {
	app := &Application{server: func(addr string, handler http.Handler) error {
		return errors.New("custom-error-from-server")
	}}

	result := app.Server("8080", nil)
	assert.Equal(t, "custom-error-from-server", result.Error())
}

func TestAppInit(t *testing.T) {
	app := &Application{init: func() error {
		return errors.New("custom-init-error")
	}}

	result := app.Initialize()
	assert.Equal(t, "custom-init-error", result.Error())
}

func TestAppRun(t *testing.T) {
	app := &Application{runner: func(args []string) error {
		return errors.New("custom-runner-error")
	}}

	result := app.Run([]string{})
	assert.Equal(t, "custom-runner-error", result.Error())
}
