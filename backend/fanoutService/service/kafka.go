package service

import (
	"context"
	"fmt"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/constants"
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/internal"
	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func PublishMessage(ctx context.Context, topicName string, message []byte) {
	tp, tperr := internal.TracerProvider(helpers.GetConfigValue("otel.endpoint"))
	if tperr != nil {
		log.Fatalf("Failed to create tracer provider: %v", tperr)
	}

	ctxNew, span := tp.Tracer("fanoutService.service").Start(ctx, "PublishMessage")
	defer span.End()

	propagators := propagation.TraceContext{}
	producer := newAccessLogProducer([]string{helpers.GetConfigValue("kafka.bootstrap.servers")}, topicName, otel.GetTracerProvider(), propagators)

	msg := sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder(topicName),
		Value: sarama.ByteEncoder(message),
	}

	propagators.Inject(ctxNew, otelsarama.NewProducerMessageCarrier(&msg))
	producer.Input() <- &msg
	successMsg := <-producer.Successes()
	fmt.Println("success, offset:", successMsg.Offset)
	span.SetAttributes(attribute.String("message", string(message)))

	err := producer.Close()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		log.Fatalf("Failed to close producer: %v", err)
	}
}

func newAccessLogProducer(brokerList []string, topicName string, tracerProvider trace.TracerProvider,
	propagators propagation.TraceContext) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}

	producer = otelsarama.WrapAsyncProducer(config, producer, otelsarama.WithTracerProvider(tracerProvider), otelsarama.WithPropagators(propagators))
	log.Println("propagators: ", propagators)

	return producer
}

func InitializeTopics() {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	admin, err := sarama.NewClusterAdmin([]string{helpers.GetConfigValue("kafka.bootstrap.servers")}, config)
	if err != nil {
		log.Printf("Failed to create cluster admin: %v", err)
	}

	admin.CreateTopic(constants.CreateTweetTopic, &sarama.TopicDetail{
		NumPartitions:     5,
		ReplicationFactor: 1,
	}, false)
	admin.CreateTopic(constants.FollowUserTopic, &sarama.TopicDetail{
		NumPartitions:     5,
		ReplicationFactor: 1,
	}, false)
	admin.CreateTopic(constants.UnfollowUserTopic, &sarama.TopicDetail{
		NumPartitions:     5,
		ReplicationFactor: 1,
	}, false)
}
