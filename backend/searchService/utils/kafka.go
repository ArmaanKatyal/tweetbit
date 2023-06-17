package utils

import (
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/helpers"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaHandler struct {
	consumer  *kafka.Consumer
	serverURL string
	groupID   string
}

// NewKafkaHandler creates a new KafkaHandler
func NewKafkaHandler(bootstrapServerURL string) (*KafkaHandler, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServerURL,
		"group.id":          helpers.GetConfigValue("kafka.group.id"),
	})
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
		return nil, err
	}

	return &KafkaHandler{
		consumer:  consumer,
		serverURL: bootstrapServerURL,
		groupID:   helpers.GetConfigValue("kafka.group.id"),
	}, nil
}

// GetConsumer returns the consumer
func (k *KafkaHandler) GetConsumer() *kafka.Consumer {
	return k.consumer
}

// GetServerURL returns the server URL
func (k *KafkaHandler) GetServerURL() string {
	return k.serverURL
}

// GetGroupID returns the group ID
func (k *KafkaHandler) GetGroupID() string {
	return k.groupID
}

// Subscribe subscribes to the given topics
func (k *KafkaHandler) Subscribe(topics []string) error {
	err := k.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("Error subscribing to topics: %v", err)
		return err
	}

	return nil
}

func (k *KafkaHandler) SubscribeToOne(topic string) error {
	err := k.consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
		return err
	}
	return nil
}
