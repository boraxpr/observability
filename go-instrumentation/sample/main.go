package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	otelScopeName    = "go-sample"
	otelScopeVersion = "v1.0.0"
)

var (
	// Array of random HTTP code
	frequent_http_code = [6]int{301, 302, 400, 403, 404, 500}
)

func main() {
	ctx := context.Background()

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, _ := prometheus.New()
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(otelScopeName, api.WithInstrumentationVersion(otelScopeVersion))

	// Adding Gauge metric
	gauge, _ := meter.Float64Gauge(
		"training_queue_length",
		api.WithDescription("The number of items in the queue."),
	)
	// Set Gauge to 1.0 but in metric will show 1
	gauge.Record(ctx, 1.0)

	// Adding Counter metric
	counter, _ := meter.Int64Counter(
		"training_response_total",
		api.WithDescription("The number of response."),
	)

	// Adding Histogram metric
	histogram, _ := meter.Float64Histogram(
		"training_http_request_duration_seconds",
		api.WithDescription("A histogram of the HTTP request durations in seconds."),
		// Override the default upper bucket bounds for our latency profile.
		api.WithExplicitBucketBoundaries(0.01, 0.02, 0.04, 0.08, 0.16, 0.32, 0.64, 1.28, 2.56, 5.12),
	)

	// Serve every paths
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling " + r.URL.Path + " ...")
		// Return hello world when request root path
		text := "Hello World!"
		// Return path when request other than root path
		if r.URL.Path != "/" {
			text = r.URL.Path
		}

		// 20% chance of http code other than 200
		http_code, http_text := randomHTTPCode(text)
		w.WriteHeader(http_code)
		fmt.Fprintln(w, http_text)
		if http_code != 200 {
			log.Println("Request " + r.URL.Path + " failed with HTTP code " + strconv.Itoa(http_code))
		}

		// Increasing counter when user request this page
		counter.Add(
			ctx,
			1,
			api.WithAttributes(
				attribute.String("path", r.URL.Path),
			),
		)

		// Set histogram with random timer
		if http_code == 200 {
			histogram.Record(
				ctx,
				randomTimer(),
				api.WithAttributes(
					attribute.String("path", r.URL.Path),
				),
			)
		}
	})

	// Serve the default Prometheus metrics registry over HTTP on /metrics.
	http.Handle("/metrics", promhttp.Handler())

	// Assign application port to run
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8081"
	}

	log.Println("Starting service on port " + appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, nil))

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

// Function to random HTTP code
func randomHTTPCode(text string) (int, string) {
	http_code := 200
	http_text := text
	// You will have 20% chance to get non-200 http code
	if rand.Float64() < 0.2 {
		http_code = frequent_http_code[rand.Intn(len(frequent_http_code))]
		http_text = "HTTP Code: " + strconv.Itoa(http_code)
	}
	return http_code, http_text
}

// This function will produce random latency and return latency in seconds
func randomTimer() float64 {
	// Create a new timer.
	start := time.Now()

	// Generate a random number between 0 and 1.
	randomNumber := rand.Float64()
	// If the random number is less than 0.8, sleep for a random time between 0 and 0.05 seconds.
	if randomNumber < 0.9 {
		sleepSeconds := rand.Float64() * 0.05
		time.Sleep(time.Duration(sleepSeconds * float64(time.Second)))
	} else {
		// Otherwise, sleep for a random time between 0.001 and 3 seconds.
		sleepSeconds := rand.Float64()*(3-0.001) + 0.001
		time.Sleep(time.Duration(sleepSeconds * float64(time.Second)))
	}

	// Return the timer
	return time.Since(start).Seconds()
}
