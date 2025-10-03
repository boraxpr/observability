package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	otelScopeName    = "go-practice"
	otelScopeVersion = "v1.0.0"
)

var (
	ctx                            context.Context
	histogram                      api.Float64Histogram
	counter_response_total         api.Int64Counter
	counter_batch_job_run_total    api.Int64Counter
	counter_batch_job_failed_total api.Int64Counter
	gauge_last_batch_run           api.Int64Gauge
	gauge_last_batch_success       api.Int64Gauge
)

type demoAPI struct{}

func (a demoAPI) register(mux *http.ServeMux) {
	mux.HandleFunc("/api/foo", a.foo)
	mux.HandleFunc("/api/bar", a.bar)
	mux.Handle("/metrics", promhttp.Handler())
}

func (a demoAPI) foo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	log.Println("Handling foo...")

	// Simulate a random duration that the "foo" operation needs to be completed.
	time.Sleep(25*time.Millisecond + time.Duration(rand.Float64()*150)*time.Millisecond)

	w.Write([]byte("Handled foo"))

	histogram.Record(
		ctx,
		time.Since(start).Seconds(),
		api.WithAttributes(
			attribute.String("path", r.URL.Path),
		),
	)
	counter_response_total.Add(
		ctx,
		1,
		api.WithAttributes(
			attribute.String("path", r.URL.Path),
			attribute.String("code", "200"),
		),
	)
}

func (a demoAPI) bar(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	log.Println("Handling bar...")
	// Simulate a random duration that the "bar" operation needs to be completed.
	time.Sleep(50*time.Millisecond + time.Duration(rand.Float64()*200)*time.Millisecond)

	w.Write([]byte("Handled bar"))

	histogram.Record(
		ctx,
		time.Since(start).Seconds(),
		api.WithAttributes(
			attribute.String("path", r.URL.Path),
		),
	)
	counter_response_total.Add(
		ctx,
		1,
		api.WithAttributes(
			attribute.String("path", r.URL.Path),
			attribute.String("code", "200"),
		),
	)
}

func periodicBackgroundTask() {
	log.Println("Starting background task loop...")
	bgTicker := time.NewTicker(5 * time.Second)
	for {
		start := time.Now()

		log.Println("Performing background task...")
		// Simulate a random duration that the background task needs to be completed.
		time.Sleep(1*time.Second + time.Duration(rand.Float64()*500)*time.Millisecond)

		// Simulate the background task either succeeding or failing (with a 30% probability).
		if rand.Float64() > 0.3 {
			log.Println("Background task completed successfully.")
			gauge_last_batch_success.Record(ctx, time.Now().Unix())
		} else {
			log.Println("Background task failed.")
			counter_batch_job_failed_total.Add(ctx, 1)
		}

		histogram.Record(
			ctx,
			time.Since(start).Seconds(),
			api.WithAttributes(
				attribute.String("path", "batch"),
			),
		)
		counter_batch_job_run_total.Add(ctx, 1)
		gauge_last_batch_run.Record(ctx, time.Now().Unix())

		<-bgTicker.C
	}
}

func main() {
	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	ctx = context.Background()
	exporter, _ := prometheus.New()
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(otelScopeName, api.WithInstrumentationVersion(otelScopeVersion))

	histogram, _ = meter.Float64Histogram(
		"training_http_request_duration_seconds",
		api.WithDescription("A histogram of the HTTP request durations in seconds."),
		// Override the default upper bucket bounds for our latency profile.
		api.WithExplicitBucketBoundaries(0.01, 0.02, 0.04, 0.08, 0.16, 0.32, 0.64, 1.28, 2.56, 5.12),
	)
	counter_response_total, _ = meter.Int64Counter(
		"training_response_total",
		api.WithDescription("The number of response."),
	)
	counter_batch_job_failed_total, _ = meter.Int64Counter(
		"training_batch_job_failed_total",
		api.WithDescription("The number of failed batch job runs."),
	)
	counter_batch_job_run_total, _ = meter.Int64Counter(
		"training_batch_job_run_total",
		api.WithDescription("The number of total batch job runs."),
	)
	gauge_last_batch_run, _ = meter.Int64Gauge(
		"training_last_batch_run",
		api.WithDescription("The Unix timestamps of the last background job runs."),
	)
	gauge_last_batch_success, _ = meter.Int64Gauge(
		"training_last_batch_success",
		api.WithDescription("The Unix timestamps of the last successful background job runs."),
	)

	// Assign application port to run
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8082"
	}

	listenAddr := flag.String("web.listen-addr", ":"+appPort, "The address to listen on for web requests.")
	flag.Parse()

	go periodicBackgroundTask()

	api := &demoAPI{}
	api.register(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(*listenAddr, nil))

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}
