package helpers

import "github.com/spf13/viper"

// returns the value of the config variable
func GetConfigValue(key string) string {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		return "NO_VALUE_FOUND"
	}
	return value
}
