package services

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/models"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Add the tweet to the timeline of all the followers of the user
func HandleCreateTweet(message *kafka.Message, rdb *utils.RedisServer) {
	log.Printf("Create Tweet: %s", message.Value)

	// destructure the incoming message
	var jsonMessage models.ITweet
	err := json.Unmarshal(message.Value, &jsonMessage)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	followers, err := GetAllUserFollowers(jsonMessage.UserId, rdb.GetUserClient())
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	tweetClient := rdb.GetTweetClient()

	// add the tweet to followers' timeline
	var wg sync.WaitGroup
	for _, follower := range followers {
		wg.Add(1)
		go func(follower string) {
			defer wg.Done()
			err = tweetClient.LPush(context.Background(), follower, message.Value).Err()
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
		}(follower)
	}
}
