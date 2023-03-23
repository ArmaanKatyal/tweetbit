package main

import (
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/services"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
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

	kafkaHandler, err := utils.NewKafkaHandler(helpers.GetConfigValue("kafka.bootstrap.servers"))
	if err != nil {
		log.Fatal(err)
	}

	kafkaHandler.Subscribe(topics)

	services.HandleRequests(kafkaHandler.GetConsumer(), rdb)
}
