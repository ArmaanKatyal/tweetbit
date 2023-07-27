package internal

import (
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// RedisServer is a wrapper around redis client
type RedisServer struct {
	tweetClient *redis.Client
}

// NewRedisServer returns a new RedisServer with a redis client
func NewRedisServer() *RedisServer {
	log.Printf("Connected to RedisServer")
	return &RedisServer{
		tweetClient: newRedisClient(1),
	}
}

func newRedisClient(db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.port"),
		Password: "",
		DB:       db,
	})
}

// GetTweetClient returns the redis client
func (r *RedisServer) GetTweetClient() *redis.Client {
	return r.tweetClient
}

// Close closes the redis clients
func (r *RedisServer) Close() error {
	err := r.tweetClient.Close()
	if err != nil {
		return err
	}
	return nil
}
