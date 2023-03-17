package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
)

func HandleFollowUser(message *kafka.Message, rdb *redis.Client) {
	// destructure the incoming message
	var jsonMessage models.IFollowUser
	err := json.Unmarshal(message.Value, &jsonMessage)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	// add the user to the follower list of the user
	err = rdb.SAdd(context.Background(), jsonMessage.UserId, jsonMessage.FollowerId).Err()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	log.Printf("Follow User: %s\n", message.Value)
}

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
