package main

import (
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/services"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	rdb := utils.NewRedisServer(helpers.GetConfigValue("redis.port"))
	defer func() {
		err := rdb.Close()
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
	}()

	topics := []string{constants.CreateTweetTopic, constants.FollowUserTopic, constants.UnfollowUserTopic}
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": helpers.GetConfigValue("kafka.bootstrap.servers"),
		"group.id":          helpers.GetConfigValue("kafka.group.id"),
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
				go services.HandleCreateTweet(e, rdb)
			case constants.FollowUserTopic:
				go services.HandleFollowUser(e, rdb.GetUserClient())
			case constants.UnfollowUserTopic:
				go services.HandleUnfollowUser(e, rdb.GetUserClient())
			}
		case *kafka.Error:
			// TOOD: handle the error properly and log it
			log.Printf("Error: %v\n", e)
		}
	}
}
