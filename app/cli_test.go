package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

func TestBuildCli(t *testing.T) {
	cli := buildCli("123-x")
	assert.Equal(t, "Brizo", cli.Name)
	assert.Equal(t, "123-x", cli.Version)
	assert.NotEmpty(t, cli.Usage)
	assert.NotEmpty(t, cli.Commands)
}

func TestRegisterCommands(t *testing.T) {
	commands := registerCommands()
	assert.Equal(t, "start", commands[0].Name)
	assert.NotEmpty(t, commands[0].Usage)
}

func TestStartCommand(t *testing.T) {
	ctx := &cli.Context{}
	mocked := new(appMock)
	mocked.On("Init").Return(nil)
	mocked.On("Server").Return(nil)

	err := startCommand(ctx, mocked.Init, mocked.Server)
	assert.Nil(t, err)
	mocked.AssertExpectations(t)
}

type appMock struct {
	mock.Mock
}

func (a *appMock) Init() error {
	args := a.Called()
	return args.Error(0)
}

func (a *appMock) Server() error {
	args := a.Called()
	return args.Error(0)
}

func TestStartCommandFailsInit(t *testing.T) {
	ctx := &cli.Context{}
	mocked := new(appMock)
	mocked.On("Init").Return(errors.New("custom-init-error"))

	err := startCommand(ctx, mocked.Init, mocked.Server)
	assert.Equal(t, "custom-init-error", err.Error())
}

func TestStartCommandFailsServer(t *testing.T) {
	ctx := &cli.Context{}
	mocked := new(appMock)
	mocked.On("Init").Return(nil)
	mocked.On("Server").Return(errors.New("custom-server-error"))

	err := startCommand(ctx, mocked.Init, mocked.Server)
	assert.Equal(t, "custom-server-error", err.Error())
}

func TestBuildFlagList(t *testing.T) {
	flags := buildFlagList()
	assert.Equal(t, 0, len(flags))

	provider1 := func() []cli.Flag {
		return []cli.Flag{
			cli.StringFlag{Name: "one"},
			cli.StringFlag{Name: "two"},
		}
	}

	provider2 := func() []cli.Flag {
		return []cli.Flag{
			cli.StringFlag{Name: "three"},
		}
	}

	flags = buildFlagList(provider1, provider2)
	assert.Equal(t, "one", flags[0].GetName())
	assert.Equal(t, "two", flags[1].GetName())
	assert.Equal(t, "three", flags[2].GetName())
}
