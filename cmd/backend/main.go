package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dmartzol/hmm/cmd/backend/api"
	"github.com/dmartzol/hmm/internal/dao/postgres"
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
		Name:  "gateway",
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
				EnvVar: "PGDATABASE",
				Value:  "hmm-development",
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
		Action: newHmmServiceRun,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func newHmmServiceRun(c *cli.Context) error {
	port := c.String(flagPort)
	host := c.String(flagHost)
	structuredLogging := c.Bool(flagStructuredLogging)

	dbConfig := postgres.Config{
		Host:     c.String(flagDBHost),
		Port:     c.Int(flagDBPort),
		Name:     c.String(flagDBName),
		User:     c.String(flagDBUser),
		Password: c.String(flagDBPass),
	}
	db, err := postgres.New(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	address := host + ":" + port
	restAPI := api.NewAPI(structuredLogging, db)
	server := &http.Server{
		Addr:    address,
		Handler: restAPI,
	}

	restAPI.Logger.Infof("listening and serving on %s", address)
	return server.ListenAndServe()
}
