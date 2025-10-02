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

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// Serve Hello World! on /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling " + r.URL.Path + " ...")
		http_code, http_text := randomHTTPCode("Hello World!")
		w.WriteHeader(http_code)
		fmt.Fprintln(w, http_text)
		if http_code != 200 {
			log.Println("Request " + r.URL.Path + " failed with HTTP code " + strconv.Itoa(http_code))
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
