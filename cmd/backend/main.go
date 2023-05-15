package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/dmartzol/go-sdk/flags"
	"github.com/dmartzol/go-sdk/logger"
	"github.com/dmartzol/go-sdk/postgres"
	"github.com/dmartzol/hmm/cmd/backend/api"
	"github.com/urfave/cli"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	otelmetric "go.opentelemetry.io/otel/metric"
)

const (
	appName = "backend"
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
		sdkLogger := logger.New()
		sdkLogger.Errorf("error running app: %v", err)
		os.Exit(1)
	}

}

func newBackendServiceRun(c *cli.Context) error {
	ctx := context.Background()

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
	defer db.Close()

	// Initializes a new logger using provided configuration and options.
	loggerOpts := []logger.Option{
		logger.WithColor(),
		logger.WithEncoding(c.String(flags.LogsFormatFlag)),
	}
	sdkLogger := logger.NewWithOptions(loggerOpts...)

	restAPI := api.NewAPI(db, sdkLogger)

	host := c.String(flags.HostnameFlagName)
	port := c.String(flags.PortFlagName)
	address := host + ":" + port
	server := &http.Server{
		Addr:    address,
		Handler: restAPI,
	}

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("github.com/open-telemetry/opentelemetry-go/example/prometheus")

	opt := otelmetric.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	)

	// This is the equivalent of prometheus.NewCounterVec
	counter, err := meter.Float64Counter("foo", otelmetric.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(ctx, 5, opt)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	gauge, err := meter.Float64ObservableGauge("bar", otelmetric.WithDescription("a fun little gauge"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = meter.RegisterCallback(func(_ context.Context, o otelmetric.Observer) error {
		n := -10. + rng.Float64()*(90.) // [-10, 100)
		o.ObserveFloat64(gauge, n, opt)
		return nil
	}, gauge)
	if err != nil {
		log.Fatal(err)
	}

	// This is the equivalent of prometheus.NewHistogramVec
	histogram, err := meter.Float64Histogram("baz", otelmetric.WithDescription("a very nice histogram"))
	if err != nil {
		log.Fatal(err)
	}
	histogram.Record(ctx, 23, opt)
	histogram.Record(ctx, 7, opt)
	histogram.Record(ctx, 101, opt)
	histogram.Record(ctx, 105, opt)

	restAPI.Logger.Infof("listening and serving on %s", address)
	return server.ListenAndServe()
}
