package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/models"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Add the tweet to the timeline of all the followers of the user
func HandleCreateTweet(ctx context.Context, message []byte, rdb *utils.RedisServer) {
	log.Printf("Create Tweet: %s", message)

	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("userGraphService.services").Start(ctx, "HandleCreateTweet")
	defer span.End()

	span.SetAttributes(attribute.Key("message").String(string(message)))

	// destructure the incoming message
	var jsonMessage models.ITweet
	err := json.Unmarshal(message, &jsonMessage)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		log.Printf("Error: %v\n", err)
	}

	// get all the followers of the user from the redis cache
	followers, err := GetAllUserFollowers(ctx, jsonMessage.UserId, rdb.GetUserClient())
	if err != nil {
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
	}

	tweetClient := rdb.GetTweetClient()

	// add the tweet to followers' timeline
	for _, follower := range followers {
		go func(follower string, ctx context.Context) {
			err = tweetClient.LPush(ctx, follower, message).Err()
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				log.Printf("Error: %v\n", err)
			}
		}(follower, ctx)
	}
}
