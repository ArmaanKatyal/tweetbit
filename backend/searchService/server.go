package main

import (
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/services"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/utils"
	"github.com/elastic/go-elasticsearch/v7"
)

func main() {
	topic := constants.CreateTweetTopic
	kakfaHandler, err := utils.NewKafkaHandler(helpers.GetConfigValue("kafka.bootstrap.servers"))
	if err != nil {
		log.Fatalf("Error while creating kafka handler: %v", err)
	}

	kakfaHandler.SubscribeToOne(topic)

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error while creating elasticsearch client: %v", err)
	}

	services.HandleRequests(kakfaHandler.GetConsumer(), es)
}
