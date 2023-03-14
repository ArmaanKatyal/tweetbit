package main

import (
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/services"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	topics := []string{constants.CreateTweetTopic, constants.FollowUserTopic, constants.UnfollowUserTopic}
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": constants.KafkaServer,
		"group.id":          "user-graph-service",
	})
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// poll for new event every 100ms
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			switch *e.TopicPartition.Topic {
			case constants.CreateTweetTopic:
				go services.HandleCreateTweet(e)
			case constants.FollowUserTopic:
				go services.HandleFollowUser(e)
			case constants.UnfollowUserTopic:
				go services.HandleUnfollowUser(e)
			}
		case *kafka.Error:
			// TOOD: handle the error properly and log it
			log.Printf("Error: %v\n", e)
		}
	}
}
