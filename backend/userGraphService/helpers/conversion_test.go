package helpers

import "testing"

func TestStringToBool(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "TestStringToBool: true",
			s:    "true",
			want: true,
		},
		{
			name: "TestStringToBool: false",
			s:    "false",
			want: false,
		},
		{
			name: "TestStringToBool: empty",
			s:    "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToBool(tt.s); got != tt.want {
				t.Errorf("StringToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
