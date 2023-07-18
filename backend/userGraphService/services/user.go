package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/models"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UserController struct {
	redis   *utils.RedisServer
	metrics *internal.PromMetrics
}

// HandleFollowUser Add the follower to the user's follower list
func (uc *UserController) HandleFollowUser(ctx context.Context, message []byte) {
	start := time.Now()
	uc.metrics.IncKafkaTransaction("followUser")

	rdb := uc.redis.GetUserClient()

	log.Printf("Follow User: %s", message)

	ctxNew, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("userGraphService.services").Start(ctx, "HandleFollowUser")
	defer span.End()

	span.SetAttributes(attribute.Key("message").String(string(message)))

	// destructure the incoming message
	var jsonMessage models.IFollowUser
	err := json.Unmarshal(message, &jsonMessage)
	if err != nil {
		uc.metrics.FollowUserResponseTime.WithLabelValues(internal.Error).Observe(time.Since(start).Seconds())
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
	}

	// add the user to the follower list of the user
	err = rdb.SAdd(ctxNew, jsonMessage.UserId, jsonMessage.FollowerId).Err()
	if err != nil {
		uc.metrics.FollowUserResponseTime.WithLabelValues(internal.Error).Observe(time.Since(start).Seconds())
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
	}

	log.Printf("Follow User: %s\n", message)
	uc.metrics.FollowUserResponseTime.WithLabelValues(internal.Success).Observe(time.Since(start).Seconds())
}

// HandleUnfollowUser Unfollow a user and remove the user from the follower list of the user
func (uc *UserController) HandleUnfollowUser(ctx context.Context, message []byte) {
	start := time.Now()
	uc.metrics.IncKafkaTransaction("unfollowUser")

	rdb := uc.redis.GetUserClient()

	log.Printf("Unfollow User: %s", message)

	ctxNew, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("userGraphService.services").Start(ctx, "HandleUnfollowUser")
	defer span.End()

	span.SetAttributes(attribute.Key("message").String(string(message)))
	// destructure the incoming message
	var jsonMessage models.IFollowUser
	err := json.Unmarshal(message, &jsonMessage)
	if err != nil {
		uc.metrics.UnfollowUserResponseTime.WithLabelValues(internal.Error).Observe(time.Since(start).Seconds())
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
	}

	// remove the user from the follower list of the user
	err = rdb.SRem(ctxNew, jsonMessage.UserId, jsonMessage.FollowerId).Err()
	if err != nil {
		uc.metrics.UnfollowUserResponseTime.WithLabelValues(internal.Error).Observe(time.Since(start).Seconds())
		log.Printf("Error: %v\n", err)
		span.SetStatus(codes.Error, err.Error())
	}

	log.Printf("Unfollow User: %s\n", message)
	uc.metrics.UnfollowUserResponseTime.WithLabelValues(internal.Success).Observe(time.Since(start).Seconds())
}

// GetAllUserFollowers Get all the followers of a user from redis
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
