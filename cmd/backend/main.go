package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dmartzol/go-sdk/flags"
	"github.com/dmartzol/go-sdk/logger"
	"github.com/dmartzol/hmm/cmd/backend/api"
	"github.com/dmartzol/hmm/internal/storage/postgres"
	"github.com/urfave/cli"
)

const (
	appName = "backend"
)

var sdkLogger logger.Logger

func init() {
	sdkLogger = logger.New()
}

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
		Name:   appName,
		Action: newHmmServiceRun,
	}
	app.Flags = append(app.Flags, flags.DatabaseFlags...)
	app.Flags = append(app.Flags, flags.LoggerFlags...)
	app.Flags = append(app.Flags, flags.ServerFlags...)

	err := app.Run(os.Args)
	if err != nil {
		sdkLogger.Errorf("error running app: %v", err)
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
