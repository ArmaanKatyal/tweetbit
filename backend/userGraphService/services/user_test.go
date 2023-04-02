package services

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var miniRedisServer *miniredis.Miniredis
var userClient *redis.Client

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}

func setup() {
	miniRedisServer = mockRedis()
	userClient = redis.NewClient(&redis.Options{
		Addr: miniRedisServer.Addr(),
		DB:   0,
	})
	tweetClient = redis.NewClient(&redis.Options{
		Addr: miniRedisServer.Addr(),
		DB:   1,
	})
}

func teardown() {
	miniRedisServer.Close()
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
		err := HandleFollowUser(context.Background(), message.Value, userClient)
		assert.Nil(t, err)

		// check if the user is added to the follower list of the user
		followers, err := GetAllUserFollowers(context.Background(), "1", userClient)
		assert.Nil(t, err)
		assert.Equal(t, []string{"2"}, followers)
	})

	t.Run("should fail with an error", func(t *testing.T) {
		miniRedisServer.SetError("FAIL")
		err := HandleFollowUser(context.Background(), message.Value, userClient)
		assert.NotNil(t, err)
		assert.Equal(t, "FAIL", err.Error())
	})

	t.Run("should fail as unable to unmarshal the message", func(t *testing.T) {
		message.Value = []byte(`{"user_id": "1", "follower_id": "2"`)
		err := HandleFollowUser(context.Background(), message.Value, userClient)
		assert.NotNil(t, err)
		assert.Equal(t, "unexpected end of JSON input", err.Error())
	})
}

func TestHandleUnFollowUser(t *testing.T) {
	setup()
	defer teardown()

	// create a message
	message := &kafka.Message{
		Value: []byte(`{"user_id": "1", "follower_id": "2"}`),
	}

	// add a follower to the user
	userClient.SAdd(context.Background(), "1", "2")

	t.Run("should remove a follower", func(t *testing.T) {
		// call the function
		err := HandleUnfollowUser(context.Background(), message.Value, userClient)
		assert.Nil(t, err)

		// check if the user is removed from the follower list of the user
		followers, err := GetAllUserFollowers(context.Background(), "1", userClient)
		assert.Nil(t, err)
		assert.Equal(t, []string{}, followers)
	})

	t.Run("should fail with an error", func(t *testing.T) {
		miniRedisServer.SetError("FAIL")
		err := HandleUnfollowUser(context.Background(), message.Value, userClient)
		assert.NotNil(t, err)
		assert.Equal(t, "FAIL", err.Error())
	})

	t.Run("should fail as unable to unmarshal the message", func(t *testing.T) {
		message.Value = []byte(`{"user_id": "1", "follower_id": "2"`)
		err := HandleUnfollowUser(context.Background(), message.Value, userClient)
		assert.NotNil(t, err)
		assert.Equal(t, "unexpected end of JSON input", err.Error())
	})
}

func TestGetAllUserFollowers(t *testing.T) {
	setup()
	defer teardown()

	// add a follower to the user
	userClient.SAdd(context.Background(), "1", "2")

	t.Run("should return all the followers of the user", func(t *testing.T) {
		followers, err := GetAllUserFollowers(context.Background(), "1", userClient)
		assert.Nil(t, err)
		assert.Equal(t, []string{"2"}, followers)
	})

	t.Run("should fail with an error", func(t *testing.T) {
		miniRedisServer.SetError("FAIL")
		_, err := GetAllUserFollowers(context.Background(), "1", userClient)
		assert.NotNil(t, err)
		assert.Equal(t, "FAIL", err.Error())
	})
}
