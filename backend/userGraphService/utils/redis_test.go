package utils

import (
	"testing"
)

func TestNewRedisServer(t *testing.T) {
	tests := []struct {
		name string
		port string
	}{
		{
			name: "TestNewRedisServer: port",
			port: "localhost:6379",
		},
		{
			name: "TestNewRedisServer: empty port",
			port: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewRedisServer(tt.port)
		})
	}
}

func TestRedisServer_GetUserClient(t *testing.T) {
	tests := []struct {
		name string
		port string
		want string
	}{
		{
			name: "TestRedisServer_GetUserClient: port",
			port: "localhost:6379",
			want: "localhost:6379",
		},
		{
			name: "TestRedisServer_GetUserClient: empty port",
			port: "",
			want: "localhost:6379",
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
	tests := []struct {
		name string
		port string
		want string
	}{
		{
			name: "TestRedisServer_GetTweetClient: port",
			port: "localhost:6379",
			want: "localhost:6379",
		},
		{
			name: "TestRedisServer_GetTweetClient: empty port",
			port: "",
			want: "localhost:6379",
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
	tests := []struct {
		name string
		port string
	}{
		{
			name: "TestRedisServer_Close: port",
			port: "localhost:6379",
		},
		{
			name: "TestRedisServer_Close: empty port",
			port: "",
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
