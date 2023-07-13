package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	es "github.com/ArmaanKatyal/tweetbit/backend/searchService/elasticsearch"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/models"
	esv7 "github.com/elastic/go-elasticsearch/v7"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TweetController struct {
	Elastic *esv7.Client
	Metrics *internal.PromMetrics
}

func (tc *TweetController) HandleCreateTweet(ctx context.Context, message []byte) {
	start := time.Now()
	log.Printf("Create Tweet: %s", message)
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("searchService.services").Start(ctx, "HandleCreateTweet")
	defer span.End()
	span.SetAttributes(attribute.Key("message").String(string(message)))

	var jsonMessage models.ITweet
	err := json.Unmarshal(message, &jsonMessage)
	if err != nil {
		tc.postMetrics(internal.Error, "tweets", "createTweet", start)
		log.Printf("Error unmarshalling message: %s", err)
	}

	esTweet := es.NewElasticTweet(tc.Elastic)
	err = esTweet.IndexTweet(ctx, jsonMessage)
	if err != nil {
		tc.postMetrics(internal.Error, "tweets", "createTweet", start)
		log.Printf("Error indexing tweet: %s", err)
	}

	tc.postMetrics(internal.Success, "tweets", "createTweet", start)
}

func (tc *TweetController) HandleDeleteTweet(ctx context.Context, message []byte) {
	start := time.Now()
	log.Printf("Delete Tweet: %s", message)
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("searchService.services").Start(ctx, "HandleDeleteTweet")
	defer span.End()

	span.SetAttributes(attribute.Key("message").String(string(message)))

	var jsonMessage models.ITweet
	err := json.Unmarshal(message, &jsonMessage)
	if err != nil {
		tc.postMetrics(internal.Error, "tweets", "deleteTweet", start)
		log.Printf("Error unmarshalling message: %s", err)
	}

	esTweet := es.NewElasticTweet(tc.Elastic)
	err = esTweet.DeleteTweet(ctx, jsonMessage.Id)
	if err != nil {
		tc.postMetrics(internal.Error, "tweets", "deleteTweet", start)
		log.Printf("Error deleting tweet: %s", err)
	}

	tc.postMetrics(internal.Success, "tweets", "deleteTweet", start)
}

func (tc *TweetController) postMetrics(code string, index string, topic string, seconds time.Time) {
	tc.Metrics.ObserveESResponseTime(code, index, time.Since(seconds).Seconds())
	tc.Metrics.CreateTweetResponseTimeHistogram.WithLabelValues(code).Observe(time.Since(seconds).Seconds())
	tc.Metrics.IncESTransaction(code, index)
	tc.Metrics.IncKafkaTransaction(topic)
}
