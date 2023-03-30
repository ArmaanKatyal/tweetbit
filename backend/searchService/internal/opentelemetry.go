package internal

import (
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	service     = helpers.GetConfigValue("otel.service")
	environment = helpers.GetConfigValue("otel.environment")
	id          = 1
)

func TracerProvider(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, utils.WrapErrorf(err, utils.ErrorCodeUnknown, "jaeger.New")
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

	return tp, nil
}
