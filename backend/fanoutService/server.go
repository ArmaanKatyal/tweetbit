package main

import (
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/service"
)

func main() {
	service.InitializeTopics()
	var server = service.NewFanoutServer()
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
