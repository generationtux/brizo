package main

import (
	"log"
	"os"

	"github.com/generationtux/brizo/app"
	"github.com/generationtux/brizo/auth"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/database/migrations"
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

	if os.Getenv("BRIZO_DEV") == "true" {
		brizoCLI.Commands = []cli.Command{
			cli.Command{
				Name:      "migrate",
				UsageText: "run database migrations",
				Action:    dbMigrate,
				Flags: []cli.Flag{
					cli.BoolFlag{Name: "down"},
				},
			},
		}
	}

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

// dbMigrate dev helper for running database migrations
func dbMigrate(c *cli.Context) error {
	if c.Bool("down") {
		return migrations.RunDown()
	}

	return migrations.Run()
}
