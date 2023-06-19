package services

import (
	"context"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/helpers"
	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v7"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

var ElasticClient *elasticsearch.Client

type Consumer struct{}

func StartConsumer(client *elasticsearch.Client) {
	ElasticClient = client
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange

	consumer, err := sarama.NewConsumer([]string{helpers.GetConfigValue("kafka.bootstrap.servers")}, config)
	if err != nil {
		log.Printf("Error while creating consumer: %v", err)
	}

	handlePartitions(consumer, constants.CreateTweetTopic)
}

func handleRequests(msg *sarama.ConsumerMessage) {
	propagaters := propagation.TraceContext{}
	ctx := propagaters.Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(msg))
	newCtx, span := otel.Tracer("searchService").Start(ctx, "handleRequests")
	defer span.End()

	propagaters.Inject(ctx, otelsarama.NewConsumerMessageCarrier(msg))

	span.SetAttributes(attribute.String("test-consumer-key", "test-consumer-value"))

	switch string(msg.Topic) {
	case constants.CreateTweetTopic:
		go HandleCreateTweet(newCtx, msg.Value, ElasticClient)
	case constants.DeleteTweetTopic:
		go HandleDeleteTweet(newCtx, msg.Value, ElasticClient)
	}

}

func handlePartitions(consumer sarama.Consumer, topic string) {
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Printf("Error while getting partitions for topic %s: %v", topic, err)
	}

	for _, partition := range partitions {
		pc, _ := consumer.ConsumePartition(constants.CreateTweetTopic, partition, sarama.OffsetOldest)
		wrappedPc := otelsarama.WrapPartitionConsumer(pc, otelsarama.WithPropagators(propagation.TraceContext{}))
		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				handleRequests(message)
			}
		}(wrappedPc)
	}
}
