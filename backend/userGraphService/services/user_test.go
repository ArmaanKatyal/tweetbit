package services

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var redisServer *miniredis.Miniredis
var redisClient *redis.Client

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}

func setup() {
	redisServer = mockRedis()
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
}

func teardown() {
	redisServer.Close()
}

func TestHandleFollowUser(t *testing.T) {
	setup()
	defer teardown()

	// create a message
	message := &kafka.Message{
		Value: []byte(`{"user_id": "1", "follower_id": "2"}`),
	}

	t.Run("should add a new follower", func(t *testing.T) {
		// call the function
		err := HandleFollowUser(message, redisClient)
		assert.Nil(t, err)

		// check if the user is added to the follower list of the user
		followers, err := GetAllUserFollowers("1", redisClient)
		assert.Nil(t, err)
		assert.Equal(t, []string{"2"}, followers)
	})

	t.Run("should fail with an error", func(t *testing.T) {
		redisServer.SetError("FAIL")
		err := HandleFollowUser(message, redisClient)
		assert.NotNil(t, err)
		assert.Equal(t, "FAIL", err.Error())
	})

}
