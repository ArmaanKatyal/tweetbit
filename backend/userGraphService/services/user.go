package services

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func HandleFollowUser(message *kafka.Message) {
	// TODO: write the logic to add user to the follower list of the user
	log.Printf("Follow User: %s\n", message.Value)
}

func HandleUnfollowUser(message *kafka.Message) {
	// TODO: write the logic to remove user from the follower list of the user
	log.Printf("Unfollow User: %s\n", message.Value)
}
