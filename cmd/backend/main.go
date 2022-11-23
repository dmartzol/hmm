package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	flagPort              = "port"
	flagHost              = "host"
	flagDBName            = "databaseName"
	flagDBPort            = "databasePort"
	flagDBHost            = "databaseHost"
	flagDBUser            = "databaseUser"
	flagDBPass            = "databasePassword"
	flagStructuredLogging = "structuredLoggin"
)

func main() {
	app := &cli.App{
		Name:  "hmm",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   flagHost,
				EnvVar: "HOST",
				Value:  "0.0.0.0",
			},
			&cli.StringFlag{
				Name:   flagPort,
				EnvVar: "PORT",
				Value:  "1100",
			},
			&cli.BoolTFlag{
				Name:   flagStructuredLogging,
				EnvVar: "STRUCTURED_LOGGING",
			},
			&cli.StringFlag{
				Name:   flagDBName,
				EnvVar: "POSTGRES_DB",
				Value:  "hmm-development",
			},
			&cli.StringFlag{
				Name:   flagDBUser,
				EnvVar: "POSTGRES_USER",
				Value:  "user-development",
			},
			&cli.StringFlag{
				Name:   flagDBPort,
				EnvVar: "DBPORT",
				Value:  "5432",
			},
			&cli.StringFlag{
				Name:   flagDBPass,
				EnvVar: "POSTGRES_PASSWORD",
				Value:  "",
			},
			&cli.StringFlag{
				Name:   flagDBHost,
				EnvVar: "PGHOST",
				Value:  "database",
			},
		},
		Action: newBackendServiceRun,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
