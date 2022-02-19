package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var (
	tp *tracesdk.TracerProvider
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.String("version", version),
		)),
	)
	return tp, nil
}

func setupTracerProvider(s *Settings) error {
	if !s.Telemetry.Enabled {
		return nil
	}

	url := "http://localhost:14268/api/traces"
	if jaegerUrl := s.Telemetry.Jaeger.URL; jaegerUrl != "" {
		url = jaegerUrl
	}

	var err error
	tp, err = tracerProvider(url)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	return nil
}

func teardownTracerProvider(ctx context.Context) error {
	return tp.Shutdown(ctx)
}
