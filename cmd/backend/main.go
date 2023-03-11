package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dmartzol/go-sdk/flags"
	"github.com/dmartzol/go-sdk/logger"
	newPostgres "github.com/dmartzol/go-sdk/postgres"
	"github.com/dmartzol/hmm/cmd/backend/api"
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
	flagStructuredLogging = "structuredLoggin"
)

func main() {
	app := &cli.App{
		Name:   appName,
		Action: newBackendServiceRun,
	}
	app.Flags = append(app.Flags, flags.DatabaseFlags...)
	app.Flags = append(app.Flags, flags.LoggerFlags...)
	app.Flags = append(app.Flags, flags.ServerFlags...)

	err := app.Run(os.Args)
	if err != nil {
		sdkLogger.Errorf("error running app: %v", err)
	}

}

func newBackendServiceRun(c *cli.Context) error {
	host := c.String(flags.HostnameFlagName)
	port := c.String(flags.PortFlagName)

	structuredLogging := c.Bool(flagStructuredLogging)

	postgresOpts := []newPostgres.Option{
		newPostgres.WithHost(c.String(flags.DatabaseHostnameFlag)),
		newPostgres.WithDatabaseName(c.String(flags.DatabaseNameFlag)),
		newPostgres.WithCreds(
			c.String(flags.DatabaseUserFlag),
			c.String(flags.DatabasePasswordFlag),
		),
	}
	db, err := newPostgres.New(appName, postgresOpts...)
	if err != nil {
		return fmt.Errorf("unable to initialize database: %w", err)
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
