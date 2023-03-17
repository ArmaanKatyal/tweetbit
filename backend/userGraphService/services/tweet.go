package services

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
)

func HandleCreateTweet(message *kafka.Message, rdb *redis.Client) {
	// TODO: write the logic to add tweet to timeline of the user
	log.Printf("Create Tweet: %s\n", message.Value)
}
