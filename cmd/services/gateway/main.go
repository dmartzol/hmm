package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	flagDBName            = "databaseName"
	flagDBPort            = "databasePort"
	flagDBHost            = "databaseHost"
	flagDBUser            = "databaseUser"
	flagDBPass            = "databasePassword"
	flagStructuredLogging = "structuredLoggin"

	dbname = "PGDATABASE"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "host",
				EnvVar: "HOST",
				Value:  "0.0.0.0",
			},
			&cli.StringFlag{
				Name:   "port",
				EnvVar: "PORT",
				Value:  "1100",
			},
			&cli.BoolTFlag{
				Name:   flagStructuredLogging,
				EnvVar: "STRUCTURED_LOGGIN",
			},
			&cli.StringFlag{
				Name:   flagDBName,
				EnvVar: "PGDATABASE",
				Value:  "development",
			},
			&cli.StringFlag{
				Name:   flagDBUser,
				EnvVar: "PGUSER",
				Value:  "user-development",
			},
			&cli.StringFlag{
				Name:   flagDBPort,
				EnvVar: "DBPORT",
				Value:  "5432",
			},
			&cli.StringFlag{
				Name:   flagDBPass,
				EnvVar: "PGPASSWORD",
				Value:  "",
			},
			&cli.StringFlag{
				Name:   flagDBHost,
				EnvVar: "PGHOST",
				Value:  "database",
			},
		},
		Name:   "gateway",
		Usage:  "",
		Action: newGatewayServiceRun,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
