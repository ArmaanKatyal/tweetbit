package internal

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	service     = "timelineSerivce"
	environment = "prod"
	id          = 1
)

type Tracer struct {
	tracerProvider *tracesdk.TracerProvider
}

// InitTracer initializes the tracer provider
func InitTracer() *Tracer {
	return &Tracer{
		tracerProvider: newTracer("http://localhost:14268/api/traces"),
	}
}

// GetTracer returns the tracer provider
func (t *Tracer) GetTracerProvider() *tracesdk.TracerProvider {
	return t.tracerProvider
}

func newTracer(url string) *tracesdk.TracerProvider {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		log.Fatal("Could not initialize jaeger exporter:", err.Error())
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", int64(id)),
			attribute.String("version", "1.0.0"),
			attribute.String("telemetry.sdk.language", "go"),
			attribute.String("telemetry.sdk.name", "opentelemetry"),
			attribute.String("telemetry.sdk.version", otel.Version()),
		)),
	)

	return tp
}
