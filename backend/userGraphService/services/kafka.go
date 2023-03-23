package services

import (
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// HandleRequests handles the requests from the kafka topics and calls the appropriate handlers
func HandleRequests(consumer *kafka.Consumer, rdb *utils.RedisServer) {
	for {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			switch *e.TopicPartition.Topic {
			case constants.CreateTweetTopic:
				go HandleCreateTweet(e, rdb)
			case constants.FollowUserTopic:
				go HandleFollowUser(e, rdb.GetUserClient())
			case constants.UnfollowUserTopic:
				go HandleUnfollowUser(e, rdb.GetUserClient())
			}
		case *kafka.Error:
			// TOOD: handle the error properly and log it
			log.Printf("Error: %v\n", e)
		}
	}
}
