package services

import (
	"context"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/constants"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-elasticsearch/v7"
)

func HandleRequests(ctx context.Context, consumer *kafka.Consumer, client *elasticsearch.Client) {
	for {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			switch *e.TopicPartition.Topic {
			case constants.CreateTweetTopic:
				go HandleCreateTweet(ctx, e, client)
			case constants.DeleteTweetTopic:
				go HandleDeleteTweet(ctx, e, client)
			}
		case *kafka.Error:
			// TODO: handle the error properly and log it
			log.Printf("Error: %v\n", e)
		}
	}
}
