package services

import (
	"context"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

var rdb *utils.RedisServer

type Consumer struct{}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Setup has been triggered")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		handleRequests(message)
		session.MarkMessage(message, "")
	}

	return nil
}

func StartConsumer(rdbServer *utils.RedisServer) {

	rdb = rdbServer // Set global variable

	consumerGroupHandler := Consumer{}
	// Wrap instrumentation
	propagators := propagation.TraceContext{}
	handler := otelsarama.WrapConsumerGroupHandler(&consumerGroupHandler, otelsarama.WithPropagators(propagators))

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup([]string{helpers.GetConfigValue("kafka.bootstrap.servers")}, helpers.GetConfigValue("kafka.group.id"), config)
	if err != nil {
		log.Printf("Error while creating consumer: %v", err)
	}

	err = consumer.Consume(context.Background(), []string{constants.CreateTweetTopic, constants.FollowUserTopic, constants.UnfollowUserTopic}, handler)
	if err != nil {
		log.Printf("Error while consuming: %v", err)
	}
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
		go HandleCreateTweet(newCtx, msg.Value, rdb)
	case constants.FollowUserTopic:
		go HandleFollowUser(newCtx, msg.Value, rdb.GetUserClient())
	case constants.UnfollowUserTopic:
		go HandleUnfollowUser(newCtx, msg.Value, rdb.GetUserClient())
	}
}
