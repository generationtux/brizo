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
	assert.NotEmpty(t, app.version)
	assert.NotEmpty(t, app.cli)
	assert.NotEmpty(t, app.config)
	assert.NotEmpty(t, app.server)
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

func TestConfigureCLI(t *testing.T) {
	cli := ConfigureCLI()
	assert.Equal(t, "Brizo", cli.Name)
}

func TestRegisterRunCommand(t *testing.T) {
	commands := registerCommands()
	assert.Equal(t, "run", commands[0].Name)
}

func TestGetAddress(t *testing.T) {
	os.Setenv("BRIZO_PORT", "")
	assert.Equal(t, ":8080", getAddress())

	os.Setenv("BRIZO_PORT", "1234")
	assert.Equal(t, ":1234", getAddress())
}
