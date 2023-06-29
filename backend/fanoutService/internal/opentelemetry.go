package internal

import (
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/utils"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	id = 1
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
			semconv.ServiceName(viper.GetString("otel.service")),
			attribute.String("environment", viper.GetString("otel.environment")),
			attribute.Int64("ID", int64(id)),
			attribute.String("version", "1.0.0"),
			attribute.String("telemetry.sdk.language", "go"),
			attribute.String("telemetry.sdk.name", "opentelemetry"),
			attribute.String("telemetry.sdk.version", otel.Version()),
		)),
	)

	return tp, nil
}
