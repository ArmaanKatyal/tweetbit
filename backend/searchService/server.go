package main

import (
	"context"
	"log"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/services"
	"github.com/elastic/go-elasticsearch/v7"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {

	tp, err := internal.TracerProvider(helpers.GetConfigValue("otel.endpoint"))
	if err != nil {
		log.Fatalf("Error while creating tracer provider: %v", err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	// create elasticsearch client with custom url
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{helpers.GetConfigValue("elasticsearch.url")},
	})
	if err != nil {
		log.Fatalf("Error creating elasticsearch client: %s", err)
	}
	services.ElasticClient = es

	client := services.NewKafkaClient()
	go client.ConsumeMessages()

	r := services.NewRouter(ctx, es)
	r.Run(helpers.GetConfigValue("server.port"))
}
