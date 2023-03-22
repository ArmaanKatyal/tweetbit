package helpers

import "testing"

func TestGetConfigValue(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "should return value: valid key",
			key:  "redis.port",
			want: "localhost:6379",
		},
		{
			name: "should not return value: invalid key",
			key:  "INVALID_KEY",
			want: "NO_VALUE_FOUND",
		},
		{
			name: "should not return value: empty key",
			key:  "",
			want: "NO_VALUE_FOUND",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfigValue(tt.key); got != tt.want {
				t.Errorf("GetConfigValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
