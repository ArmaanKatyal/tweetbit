package services

import (
	"context"
	"testing"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var tweetClient *redis.Client

func TestHandleCreateTweet(t *testing.T) {
	setup()
	defer teardown()

	// create fake followers for the user
	userClient.SAdd(context.Background(), "1", "2", "3")
	members, err := userClient.SMembers(context.Background(), "1").Result()
	assert.Nil(t, err)
	assert.Equal(t, []string{"2", "3"}, members)

	// create a message
	message := &kafka.Message{
		Value: []byte(`{"user_id": "1", "id": "1", "content": "test"}`),
	}

	// create a redis struct to pass to the function
	var redisStruct utils.RedisServer
	redisStruct.SetTweetClient(tweetClient)
	redisStruct.SetUserClient(userClient)

	t.Run("should add the tweet to the timeline of all the followers", func(t *testing.T) {
		// call the function
		err := HandleCreateTweet(message, &redisStruct)
		assert.Nil(t, err)

		/*
			We are using time.Sleep() to wait for the goroutine to finish which adds
			the tweet to the timeline of all the followers.
		*/
		time.Sleep(500 * time.Millisecond)
		// check if the tweet is added to the timeline of all the followers
		tweet1, err := tweetClient.LRange(context.Background(), "2", 0, -1).Result()
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, len(tweet1), 1)

		tweet2, err := tweetClient.LRange(context.Background(), "3", 0, -1).Result()
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, len(tweet2), 1)
	})

	t.Run("should fail as unable to unmarshal the message", func(t *testing.T) {
		// create a message with invalid json
		message := &kafka.Message{
			Value: []byte(`{"user_id": "1", "id": "1", "content": "test"`),
		}

		// call the function
		err := HandleCreateTweet(message, &redisStruct)
		assert.NotNil(t, err)
		assert.Equal(t, "unexpected end of JSON input", err.Error())
	})

	t.Run("should fail as unable to add the tweet to the timeline of the followers", func(t *testing.T) {
		redisServer.SetError("FAIL")
		err := HandleCreateTweet(message, &redisStruct)
		assert.NotNil(t, err)
		assert.Equal(t, "FAIL", err.Error())
	})
}
