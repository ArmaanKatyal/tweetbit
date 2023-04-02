package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/models"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Add the follower to the user's follower list
func HandleFollowUser(ctx context.Context, message []byte, rdb *redis.Client) error {
	log.Printf("Follow User: %s", message)

	ctxNew, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("userGraphService.services").Start(ctx, "HandleFollowUser")
	defer span.End()

	span.SetAttributes(attribute.Key("message").String(string(message)))

	// destructure the incoming message
	var jsonMessage models.IFollowUser
	err := json.Unmarshal(message, &jsonMessage)
	if err != nil {
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// add the user to the follower list of the user
	err = rdb.SAdd(ctxNew, jsonMessage.UserId, jsonMessage.FollowerId).Err()
	if err != nil {
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Printf("Follow User: %s\n", message)
	return nil
}

// Unfollow a user and remove the user from the follower list of the user
func HandleUnfollowUser(ctx context.Context, message []byte, rdb *redis.Client) error {
	log.Printf("Unfollow User: %s", message)

	ctxNew, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("userGraphService.services").Start(ctx, "HandleUnfollowUser")
	defer span.End()

	span.SetAttributes(attribute.Key("message").String(string(message)))
	// destructure the incoming message
	var jsonMessage models.IFollowUser
	err := json.Unmarshal(message, &jsonMessage)
	if err != nil {
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// remove the user from the follower list of the user
	err = rdb.SRem(ctxNew, jsonMessage.UserId, jsonMessage.FollowerId).Err()
	if err != nil {
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Printf("Unfollow User: %s\n", message)
	return nil
}

// Get all the followers of a user from redis
func GetAllUserFollowers(ctx context.Context, userId string, rdb *redis.Client) ([]string, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("userGraphService.services").Start(ctx, "GetAllUserFollowers")
	defer span.End()

	span.SetAttributes(attribute.Key("userId").String(string(userId)))
	followers, err := rdb.SMembers(ctx, userId).Result()
	if err != nil {
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return followers, nil
}
