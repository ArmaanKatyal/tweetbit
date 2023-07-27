package utils

import (
	"os"

	"github.com/spf13/viper"
)

// LoadConfig loads the config file
func LoadConfig(envFlag string) {
	viper.SetConfigName(envFlag)
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// returns the value of the environment variable
func GetEnvValue(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return "test"
}
