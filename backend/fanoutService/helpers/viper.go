package helpers

import "github.com/spf13/viper"

// returns the value of the environment variable
func GetEnvValue(key string) string {
	viper.SetConfigName(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		panic("Invalid type assertion")
	}
	return value
}

// returns the value of the config variable
func GetConfigValue(key string) string {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		panic("Invalid type assertion")
	}
	return value
}