package utils

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

var redisServer *miniredis.Miniredis

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}

func setup() {
	redisServer = mockRedis()
}

func teardown() {
	redisServer.Close()
}

func TestNewRedisServer(t *testing.T) {
	setup()
	defer teardown()
	tests := []struct {
		name string
		port string
	}{
		{
			name: "TestNewRedisServer: port",
			port: redisServer.Addr(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewRedisServer(tt.port)
		})
	}
}

func TestRedisServer_GetUserClient(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name string
		port string
		want string
	}{
		{
			name: "TestRedisServer_GetUserClient: port",
			port: redisServer.Addr(),
			want: redisServer.Addr(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedisServer(tt.port)
			if r.GetUserClient().Options().Addr != tt.want {
				t.Errorf("RedisServer.GetUserClient() = %v, want %v", r.GetUserClient().Options().Addr, tt.want)
			}
		})
	}
}

func TestRedisServer_GetTweetClient(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name string
		port string
		want string
	}{
		{
			name: "TestRedisServer_GetTweetClient: port",
			port: redisServer.Addr(),
			want: redisServer.Addr(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedisServer(tt.port)
			if r.GetTweetClient().Options().Addr != tt.want {
				t.Errorf("RedisServer.GetTweetClient() = %v, want %v", r.GetTweetClient().Options().Addr, tt.want)
			}
		})
	}
}

func TestRedisServer_Close(t *testing.T) {
	setup()
	defer teardown()
	tests := []struct {
		name string
		port string
	}{
		{
			name: "TestRedisServer_Close: port",
			port: redisServer.Addr(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedisServer(tt.port)
			if err := r.Close(); err != nil {
				t.Errorf("RedisServer.Close() error = %v", err)
			}
		})
	}
}

func TestRedisServer_SetUserClient(t *testing.T) {
	setup()
	defer teardown()

	t.Run("TestRedisServer_SetUserClient: port", func(t *testing.T) {
		r := NewRedisServer(redisServer.Addr())
		r.SetUserClient(r.GetUserClient())
		assert.Equal(t, r.GetUserClient().Options().Addr, redisServer.Addr())
	})
}

func TestRedisServer_SetTweetClient(t *testing.T) {
	setup()
	defer teardown()

	t.Run("TestRedisServer_SetTweetClient: port", func(t *testing.T) {
		r := NewRedisServer(redisServer.Addr())
		r.SetTweetClient(r.GetTweetClient())
		assert.Equal(t, r.GetTweetClient().Options().Addr, redisServer.Addr())
	})
}
