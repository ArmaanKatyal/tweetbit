package main

import (
	"fmt"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/service"
)

func main() {
	var server = service.NewFanoutServer()
	fmt.Println(helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableUserGraph")))
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
