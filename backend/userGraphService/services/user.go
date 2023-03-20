package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
)

// Add the follower to the user's follower list
func HandleFollowUser(message *kafka.Message, rdb *redis.Client) error {
	// destructure the incoming message
	var jsonMessage models.IFollowUser
	err := json.Unmarshal(message.Value, &jsonMessage)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}

	// add the user to the follower list of the user
	err = rdb.SAdd(context.Background(), jsonMessage.UserId, jsonMessage.FollowerId).Err()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}

	log.Printf("Follow User: %s\n", message.Value)
	return nil
}

// Unfollow a user and remove the user from the follower list of the user
func HandleUnfollowUser(message *kafka.Message, rdb *redis.Client) {
	// destructure the incoming message
	var jsonMessage models.IFollowUser
	err := json.Unmarshal(message.Value, &jsonMessage)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	// remove the user from the follower list of the user
	err = rdb.SRem(context.Background(), jsonMessage.UserId, jsonMessage.FollowerId).Err()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	log.Printf("Unfollow User: %s\n", message.Value)
}

// Get all the followers of a user from redis
func GetAllUserFollowers(userId string, rdb *redis.Client) ([]string, error) {
	followers, err := rdb.SMembers(context.Background(), userId).Result()
	if err != nil {
		return nil, err
	}
	return followers, nil
}
