package main

import (
	"context"
	"log"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/services"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	tp, err := internal.TracerProvider(helpers.GetConfigValue("otel.endpoint"))
	if err != nil {
		log.Printf("Error while creating tracer provider: %v", err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error while shutting down tracer provider: %v", err)
		}
	}(ctx)

	rdb := utils.NewRedisServer(helpers.GetConfigValue("redis.port"))
	defer func() {
		err := rdb.Close()
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
	}()
	services.RDB = rdb

	client := services.NewKafkaClient()
	go client.ConsumeMessages()

	r := services.NewRouter()
	r.Run(helpers.GetConfigValue("server.port"))
}
