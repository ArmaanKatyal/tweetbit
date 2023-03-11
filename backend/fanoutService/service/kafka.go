package service

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Kafka struct {
	producer      *kafka.Producer
	topic         string
	delivery_chan chan kafka.Event
}

func NewKafka(p *kafka.Producer, t string) *Kafka {
	return &Kafka{
		producer:      p,
		topic:         t,
		delivery_chan: make(chan kafka.Event),
	}
}

// publish a message to kafka topic
func (op *Kafka) PublishMessage(message []byte) error {
	err := op.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &op.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, op.delivery_chan)
	if err != nil {
		return err
	}

	e := <-op.delivery_chan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	return nil
}
