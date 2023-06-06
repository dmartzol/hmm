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
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dmartzol/hmm/cmd/backend/api"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

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

	// Initializes a new prometheus exporter and registers it as a metrics provider.
	promExporter, err := prometheus.New()
	if err != nil {
		restAPI.Logger.Errorf("failed to initialize prometheus exporter: %v", err)
		os.Exit(1)
	}
	provider := metric.NewMeterProvider(metric.WithReader(promExporter))
	meter := provider.Meter("github.com/open-telemetry/opentelemetry-go/example/prometheus")
	sampleMetrics(ctx, meter)

	err = initTraces(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize trace collector and exporter: %w", err)
	}

	restAPI.Logger.Infof("listening and serving on %s", address)
	return server.ListenAndServe()
}

func sampleMetrics(ctx context.Context, meter otelmetric.Meter) {
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
}

func initTraces(ctx context.Context) error {
	// Initializes a new grpc connection to the collector.
	conn, err := grpc.DialContext(ctx, "otel:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial collector: %w", err)
	}

	// Set up a trace exporter
	otlpOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithEndpoint("tempo:4317"),
		otlptracegrpc.WithInsecure(),
	}
	traceExporter, err := otlptracegrpc.New(ctx, otlpOpts...)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Set up a trace provider
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(tracerProvider)

	return nil
}
