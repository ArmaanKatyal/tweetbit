package services

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func HandleCreateTweet(message *kafka.Message) {
	// TODO: write the logic to add tweet to timeline of the user
	log.Printf("Create Tweet: %s\n", message.Value)
}
