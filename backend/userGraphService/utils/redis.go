package utils

import "github.com/redis/go-redis/v9"

// RedisServer is a wrapper around redis client
type RedisServer struct {
	client *redis.Client
}

// NewRedisServer returns a new RedisServer with a redis client
func NewRedisServer() *RedisServer {
	return &RedisServer{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

// GetClient returns the redis client
func (r *RedisServer) GetClient() *redis.Client {
	return r.client
}

// Close closes the redis client
func (r *RedisServer) Close() error {
	return r.client.Close()
}
