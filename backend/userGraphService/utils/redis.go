package utils

import (
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisServer is a wrapper around redis client
type RedisServer struct {
	userClient  *redis.Client
	tweetClient *redis.Client
}

// NewRedisServer returns a new RedisServer with a redis client
func NewRedisServer(port string) *RedisServer {
	log.Printf("Connected to RedisServer")
	return &RedisServer{
		userClient: redis.NewClient(&redis.Options{
			Addr:     port,
			Password: "",
			DB:       0,
		}),
		tweetClient: redis.NewClient(&redis.Options{
			Addr:     port,
			Password: "",
			DB:       1,
		}),
	}
}

// GetUserClient returns the redis client
func (r *RedisServer) GetUserClient() *redis.Client {
	return r.userClient
}

// GetTweetClient returns the redis client
func (r *RedisServer) GetTweetClient() *redis.Client {
	return r.tweetClient
}

// Close closes the redis clients
func (r *RedisServer) Close() error {
	err := r.userClient.Close()
	if err != nil {
		return err
	}
	err = r.tweetClient.Close()
	if err != nil {
		return err
	}
	return nil
}
