package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
