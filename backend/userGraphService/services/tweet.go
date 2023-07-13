package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/models"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type TweetController struct {
	redis   *utils.RedisServer
	metrics *internal.PromMetrics
}

// Add the tweet to the timeline of all the followers of the user
func (tc *TweetController) HandleCreateTweet(ctx context.Context, message []byte) {
	start := time.Now()
	tc.metrics.IncKafkaTransaction("createTweet")

	log.Printf("Create Tweet: %s", message)

	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("userGraphService.services").Start(ctx, "HandleCreateTweet")
	defer span.End()

	span.SetAttributes(attribute.Key("message").String(string(message)))

	// destructure the incoming message
	var jsonMessage models.ITweet
	err := json.Unmarshal(message, &jsonMessage)
	if err != nil {
		tc.metrics.CreateTweetResponseTime.WithLabelValues(internal.Error).Observe(time.Since(start).Seconds())
		span.SetStatus(codes.Error, err.Error())
		log.Printf("Error: %v\n", err)
	}

	// get all the followers of the user from the redis cache
	followers, err := GetAllUserFollowers(ctx, jsonMessage.UserId, tc.redis.GetUserClient())
	if err != nil {
		tc.metrics.CreateTweetResponseTime.WithLabelValues(internal.Error).Observe(time.Since(start).Seconds())
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
	}

	tweetClient := tc.redis.GetTweetClient()

	// add the tweet to followers' timeline
	for _, follower := range followers {
		// TODO: fix this go routine logic. instead of starting a new go routine for each follower, use a channel
		go func(follower string, ctx context.Context) {
			err = tweetClient.LPush(ctx, follower, message).Err()
			if err != nil {
				tc.metrics.CreateTweetResponseTime.WithLabelValues(internal.Error).Observe(time.Since(start).Seconds())
				span.SetStatus(codes.Error, err.Error())
				log.Printf("Error: %v\n", err)
			}
		}(follower, ctx)
	}
	tc.metrics.CreateTweetResponseTime.WithLabelValues(internal.Success).Observe(time.Since(start).Seconds())
}
