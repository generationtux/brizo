package main

import (
	"log"
	"os"

	"github.com/generationtux/brizo/app"
	"github.com/generationtux/brizo/auth"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {
	// .env file for local configuration during development (see .env.example)
	godotenv.Load()

	version := "0.1.0" // @todo dynamically read current version

	brizoCLI := cli.NewApp()
	brizoCLI.Name = "Brizo"
	brizoCLI.Usage = "PAAS powered by Kubernetes."
	brizoCLI.UsageText = "brizo [options]"
	brizoCLI.Version = version
	brizoCLI.Action = runCLI
	brizoCLI.Flags = buildFlagList(
		auth.CLIFlags(),
		kube.CLIFlags(),
		database.CLIFlags(),
	)

	brizoCLI.Run(os.Args)
}

// runCLI initializes Brizo and starts the HTTP server
func runCLI(c *cli.Context) error {
	brizo := app.New()

	log.Println("Intializing app...")
	err := brizo.Initialize()
	if err != nil {
		return err
	}

	log.Println("Starting server...")
	return brizo.Server()
}

// buildFlags pulls together configuration flags
func buildFlagList(flagGroups ...[]cli.Flag) []cli.Flag {
	var flags []cli.Flag

	for _, group := range flagGroups {
		flags = append(flags, group...)
	}

	return flags
}
