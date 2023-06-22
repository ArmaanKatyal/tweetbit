package utils

// import (
// 	"testing"

// 	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/helpers"
// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// 	"github.com/stretchr/testify/assert"
// )

// func TestNewKafkaHandler(t *testing.T) {
// 	mock, err := kafka.NewMockCluster(1)
// 	assert.Nil(t, err)
// 	defer mock.Close()
// 	t.Run("should create a new kafka handler", func(t *testing.T) {
// 		handler, err := NewKafkaHandler(mock.BootstrapServers())
// 		assert.Nil(t, err)

// 		assert.Equal(t, handler.GetServerURL(), mock.BootstrapServers())
// 	})
// }

// func TestKafkaHandler_GetConsumer(t *testing.T) {
// 	mock, err := kafka.NewMockCluster(1)
// 	assert.Nil(t, err)
// 	defer mock.Close()
// 	t.Run("should return the consumer", func(t *testing.T) {
// 		handler, err := NewKafkaHandler(mock.BootstrapServers())
// 		assert.Nil(t, err)

// 		assert.Equal(t, handler.GetConsumer(), handler.consumer)
// 	})
// }

// func TestKafkaHandler_GetServerURL(t *testing.T) {
// 	mock, err := kafka.NewMockCluster(1)
// 	assert.Nil(t, err)
// 	defer mock.Close()
// 	t.Run("should return the server URL", func(t *testing.T) {
// 		handler, err := NewKafkaHandler(mock.BootstrapServers())
// 		assert.Nil(t, err)

// 		assert.Equal(t, handler.GetServerURL(), mock.BootstrapServers())
// 	})
// }

// func TestKafkaHandler_GetGroupID(t *testing.T) {
// 	mock, err := kafka.NewMockCluster(1)
// 	assert.Nil(t, err)
// 	defer mock.Close()
// 	t.Run("should return the group ID", func(t *testing.T) {
// 		handler, err := NewKafkaHandler(mock.BootstrapServers())
// 		assert.Nil(t, err)

// 		assert.Equal(t, handler.GetGroupID(), helpers.GetConfigValue("kafka.group.id"))
// 	})
// }

// func TestKafkaHandler_Subscribe(t *testing.T) {
// 	mock, err := kafka.NewMockCluster(1)
// 	assert.Nil(t, err)
// 	defer mock.Close()
// 	t.Run("should subscribe to the given topics", func(t *testing.T) {
// 		handler, err := NewKafkaHandler(mock.BootstrapServers())
// 		assert.Nil(t, err)

// 		err = handler.Subscribe([]string{"test"})
// 		assert.Nil(t, err)
// 	})
// }
