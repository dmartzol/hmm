package main

import (
	"log"
	"os"

	"github.com/dmartzol/hmm/internal/service"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:   "gateway",
		Usage:  "",
		Action: service.NewGatewayServiceRun,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
