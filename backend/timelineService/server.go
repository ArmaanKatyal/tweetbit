package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/services"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// parse the environment flag to determine which config file to load
var envFlag string

func init() {
	flag.StringVar(&envFlag, "env", "prod", "environment flag")
	flag.Parse()
}

func main() {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if envFlag == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	utils.LoadConfig(envFlag)
	pm := internal.InitPromMetrics("timelineservice", prometheus.LinearBuckets(0, 5, 20))
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
	r.Run(viper.GetString("server.port"))
}
