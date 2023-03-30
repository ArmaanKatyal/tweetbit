package services

import (
	"context"
	"encoding/json"
	"log"

	es "github.com/ArmaanKatyal/tweetbit/backend/searchService/elasticsearch"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/models"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/utils"
	esv7 "github.com/elastic/go-elasticsearch/v7"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func HandleCreateTweet(ctx context.Context, message *kafka.Message, client *esv7.Client) error {
	log.Printf("Create Tweet: %s", message.Value)

	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("searchService.services").Start(ctx, "HandleCreateTweet")
	defer span.End()
	span.SetAttributes(attribute.Key("message").String(string(message.Value)))

	var jsonMessage models.ITweet
	err := json.Unmarshal(message.Value, &jsonMessage)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "json.Unmarshal")
	}

	esTweet := es.NewElasticTweet(client)
	err = esTweet.IndexTweet(ctx, jsonMessage)
	if err != nil {
		log.Printf("Error indexing tweet: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "esTweet.IndexTweet")
	}

	return nil
}

func HandleDeleteTweet(ctx context.Context, message *kafka.Message, client *esv7.Client) error {
	log.Printf("Delete Tweet: %s", message.Value)

	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("searchService.services").Start(ctx, "HandleDeleteTweet")
	defer span.End()

	span.SetAttributes(attribute.Key("message").String(string(message.Value)))

	var jsonMessage models.ITweet
	err := json.Unmarshal(message.Value, &jsonMessage)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "json.Unmarshal")
	}

	esTweet := es.NewElasticTweet(client)
	err = esTweet.DeleteTweet(ctx, jsonMessage.Id)
	if err != nil {
		log.Printf("Error deleting tweet: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "esTweet.DeleteTweet")
	}

	return nil
}
