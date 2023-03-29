package services

import (
	"encoding/json"
	"log"

	es "github.com/ArmaanKatyal/tweetbit/backend/searchService/elasticsearch"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/models"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/utils"
	esv7 "github.com/elastic/go-elasticsearch/v7"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func HandleCreateTweet(message *kafka.Message, client *esv7.Client) error {
	log.Printf("Create Tweet: %s", message.Value)

	var jsonMessage models.ITweet
	err := json.Unmarshal(message.Value, &jsonMessage)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "json.Unmarshal")
	}

	esTweet := es.NewElasticTweet(client)
	err = esTweet.IndexTweet(jsonMessage)
	if err != nil {
		log.Printf("Error indexing tweet: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "esTweet.IndexTweet")
	}

	return nil
}

func HandleDeleteTweet(message *kafka.Message, client *esv7.Client) error {
	log.Printf("Delete Tweet: %s", message.Value)

	var jsonMessage models.ITweet
	err := json.Unmarshal(message.Value, &jsonMessage)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "json.Unmarshal")
	}

	esTweet := es.NewElasticTweet(client)
	err = esTweet.DeleteTweet(jsonMessage.Id)
	if err != nil {
		log.Printf("Error deleting tweet: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "esTweet.DeleteTweet")
	}

	return nil
}
