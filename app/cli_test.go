package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCli(t *testing.T) {

}

func TestConfigureCLI(t *testing.T) {
	cli := configureCLI()
	assert.Equal(t, "Brizo", cli.Name)
}

func TestRegisterRunCommand(t *testing.T) {
	commands := registerCommands()
	assert.Equal(t, "run", commands[0].Name)
}
