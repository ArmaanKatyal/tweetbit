package main

import (
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
)

func main() {
	var server = helpers.NewFanoutServer()
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
