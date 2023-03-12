package test

import (
	"testing"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
)

func TestStringToBool(t *testing.T) {
	// create a test table
	var testTable = []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"something", false},
	}

	for _, test := range testTable {
		if helpers.StringToBool(test.input) != test.expected {
			t.Errorf("Test failed")
		}
	}

}

func TestGetConfigValue(t *testing.T) {
	// create a test table
	var testTable = []struct {
		input    string
		expected string
	}{
		{"server.host", "localhost"},
		{"something", "NO_VALUE_FOUND"},
	}

	for _, test := range testTable {
		if helpers.GetConfigValue(test.input) != test.expected {
			t.Errorf("Test failed")
		}
	}

}
