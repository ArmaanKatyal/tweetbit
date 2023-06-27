package main

import (
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/service"
	"github.com/spf13/viper"
)

func main() {
	helpers.LoadConfig()
	if viper.GetBool("featureFlag.enableTopicCreation") {
		service.InitializeTopics()
	}
	var server = service.NewFanoutServer()
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
