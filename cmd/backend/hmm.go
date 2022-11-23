package main

import (
	"log"
	"net/http"

	"github.com/dmartzol/hmm/cmd/backend/api"
	"github.com/dmartzol/hmm/internal/storage/postgres"
	"github.com/urfave/cli"
)

func newBackendServiceRun(c *cli.Context) error {
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
