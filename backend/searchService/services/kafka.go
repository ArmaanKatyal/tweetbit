package services

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func HandleRequests(consumer *kafka.Consumer) {
	for {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			log.Printf("Message on %s: %s\n", e.TopicPartition, string(e.Value))
			// do something
		case *kafka.Error:
			// TODO: handle the error properly and log it
			log.Printf("Error: %v\n", e)
		}
	}
}
