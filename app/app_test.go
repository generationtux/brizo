package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestNewApplication(t *testing.T) {
	app := New()
	assert.NotEmpty(t, app.version)
	assert.NotEmpty(t, app.server)
	assert.NotEmpty(t, app.builder)
}

func TestAppVersion(t *testing.T) {
	app := &Application{version: "1.x"}
	assert.Equal(t, "1.x", app.Version())
}

func TestAppServer(t *testing.T) {
	app := &Application{server: func() error {
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

func TestAppCliBuilder(t *testing.T) {
	app := &Application{builder: func(version string) *cli.App {
		return &cli.App{Name: "custom"}
	}}

	result := app.BuildCLI()
	assert.Equal(t, "custom", result.Name)
}
