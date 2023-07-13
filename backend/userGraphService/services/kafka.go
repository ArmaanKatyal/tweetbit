package services

import (
	"context"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/Shopify/sarama"
)

// var RDB *utils.RedisServer

type KafkaClient struct {
	consumer sarama.ConsumerGroup
	redis    *utils.RedisServer
	metrics  *internal.PromMetrics
}

func NewKafkaClient(rdb *utils.RedisServer, pm *internal.PromMetrics) *KafkaClient {
	return &KafkaClient{
		consumer: createKafkaConsumer([]string{helpers.GetConfigValue("kafka.bootstrap.servers")}),
		redis:    rdb,
		metrics:  pm,
	}
}

func createKafkaConsumer(brokers []string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange

	consumer, err := sarama.NewConsumerGroup(brokers, helpers.GetConfigValue("kafka.group"), config)
	if err != nil {
		log.Printf("Error while creating consumer: %v", err)
	}

	return consumer
}

func (kc *KafkaClient) ConsumeMessages() {
	consumerGroupHandler := Consumer{}
	consumerGroupHandler.Redis = kc.redis
	consumerGroupHandler.Metrics = kc.metrics

	// Wrap instrumentation
	propagators := propagation.TraceContext{}
	handler := otelsarama.WrapConsumerGroupHandler(&consumerGroupHandler, otelsarama.WithPropagators(propagators))

	err := kc.consumer.Consume(context.Background(), []string{constants.CreateTweetTopic, constants.FollowUserTopic, constants.UnfollowUserTopic}, handler)
	if err != nil {
		log.Printf("Error while consuming messages: %v", err)
	}
}

func handleRequests(msg *sarama.ConsumerMessage, rdb *utils.RedisServer, pm *internal.PromMetrics) {
	propagaters := propagation.TraceContext{}
	ctx := propagaters.Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(msg))
	newCtx, span := otel.Tracer("searchService").Start(ctx, "handleRequests")
	defer span.End()

	propagaters.Inject(ctx, otelsarama.NewConsumerMessageCarrier(msg))

	tc := &TweetController{
		redis:   rdb,
		metrics: pm,
	}
	uc := &UserController{
		redis:   rdb,
		metrics: pm,
	}

	switch string(msg.Topic) {
	case constants.CreateTweetTopic:
		go tc.HandleCreateTweet(newCtx, msg.Value)
	case constants.FollowUserTopic:
		go uc.HandleFollowUser(newCtx, msg.Value)
	case constants.UnfollowUserTopic:
		go uc.HandleUnfollowUser(newCtx, msg.Value)
	}
}

type Consumer struct {
	Redis   *utils.RedisServer
	Metrics *internal.PromMetrics
}

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
	// for message := range claim.Messages() {
	// 	handleRequests(message)
	// 	session.MarkMessage(message, "")
	// }
	for message := range claim.Messages() {
		handleRequests(message, consumer.Redis, consumer.Metrics)
		session.MarkMessage(message, "")
	}

	return nil
}
