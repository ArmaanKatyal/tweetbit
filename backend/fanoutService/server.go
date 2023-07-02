package main

import (
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

func main() {
	helpers.LoadConfig()
	if viper.GetBool("featureFlag.enableTopicCreation") {
		service.InitializeTopics()
	}

	pm := internal.InitPromMetrics("fanoutservice", prometheus.LinearBuckets(0, 5, 20))

	var server = service.NewFanoutServer(pm)
	go server.Run()

	r := service.NewRouter(pm)
	r.Run(viper.GetString("gin.port"))
}
