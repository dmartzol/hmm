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
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

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

	shutdown, err := initProvider(ctx, appName)
	if err != nil {
		return fmt.Errorf("failed to initialize trace provider: %w", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

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

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "otel:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}
