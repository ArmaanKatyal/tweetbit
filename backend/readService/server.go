package main

import (
	"context"
	"log"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	pm := internal.InitPromMetrics("readservice", prometheus.LinearBuckets(0, 5, 20))
	tp := internal.InitTracer()

	otel.SetTracerProvider(tp.GetTracerProvider())
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutting down
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := tp.GetTracerProvider().Shutdown(ctx); err != nil {
			log.Fatal("Failed to shutdown tracer provider: ", err.Error())
		}
	}(ctx)

	// start the server
	r := services.NewRouter(ctx, pm)
	r.Run(":5005")
}
