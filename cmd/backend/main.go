package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dmartzol/go-sdk/flags"
	"github.com/dmartzol/go-sdk/logger"
	"github.com/dmartzol/go-sdk/postgres"
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
	postgresOpts := []postgres.Option{
		postgres.WithHost(c.String(flags.DatabaseHostnameFlag)),
		postgres.WithDatabaseName(c.String(flags.DatabaseNameFlag)),
		postgres.WithCreds(
			c.String(flags.DatabaseUserFlag),
			c.String(flags.DatabasePasswordFlag),
		),
	}
	db, err := postgres.New(appName, postgresOpts...)
	if err != nil {
		return fmt.Errorf("unable to initialize database: %w", err)
	}

	// Initializes a new logger using provided configuration and options.
	loggerOpts := []logger.Option{
		logger.WithColor(),
		logger.WithEncoding(c.String(flags.LogsFormatFlag)),
	}
	sdkLogger = logger.NewWithOptions(loggerOpts...)

	restAPI := api.NewAPI(db, sdkLogger)

	host := c.String(flags.HostnameFlagName)
	port := c.String(flags.PortFlagName)
	address := host + ":" + port
	server := &http.Server{
		Addr:    address,
		Handler: restAPI,
	}

	restAPI.Logger.Infof("listening and serving on %s", address)
	return server.ListenAndServe()
}
