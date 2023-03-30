package main

import (
	"context"
	"log"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/services"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/elastic/go-elasticsearch/v7"
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

	ctx, span := tp.Tracer("searchService").Start(ctx, "searchService.main")
	defer span.End()

	topic := constants.CreateTweetTopic
	kakfaHandler, err := utils.NewKafkaHandler(helpers.GetConfigValue("kafka.bootstrap.servers"))
	if err != nil {
		log.Fatalf("Error while creating kafka handler: %v", err)
	}

	kakfaHandler.SubscribeToOne(topic)

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error while creating elasticsearch client: %v", err)
	}

	services.HandleRequests(ctx, kakfaHandler.GetConsumer(), es)
}
