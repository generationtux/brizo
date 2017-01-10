package app

import (
	"github.com/generationtux/brizo/auth"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/urfave/cli"
)

// buildCli sets up the CLI app
func buildCli(version string) *cli.App {
	brizoCLI := cli.NewApp()
	brizoCLI.Name = "Brizo"
	brizoCLI.Usage = "PAAS powered by Kubernetes."
	brizoCLI.Version = version
	brizoCLI.Commands = registerCommands()

	return brizoCLI
}

// registerCommands sets up the CLI commands
func registerCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "start",
			Usage: "start Brizo",
			Action: func(c *cli.Context) error {
				return nil
			},
			Flags: buildFlagList(auth.CLIFlags, database.CLIFlags, kube.CLIFlags),
		},
	}
}

// startCommand initializes Brizo and starts the server
func startCommand(c *cli.Context, init Initializer, server HTTPServer) error {
	err := init()
	if err != nil {
		return err
	}

	return server()
}

// buildFlagList joins the provided flag sources into one slice of cli.Flags
func buildFlagList(providers ...FlagProvider) []cli.Flag {
	var flags []cli.Flag
	for _, provider := range providers {
		flags = append(flags, provider()...)
	}

	return flags
}
