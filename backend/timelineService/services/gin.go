package services

import (
	"context"

	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/controllers"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewRouter(ctx context.Context, pm *internal.PromMetrics) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	ds := internal.NewDatabaseService(&internal.DatabaseConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
	})

	health := &controllers.HealthController{}
	health.Metrics = pm

	router.GET("/health", health.Status)
	router.GET("/metrics", internal.PrometheusHandler())
	router.Use(middleware.VerifyApiKey(pm))

	v1 := router.Group("/api/v1")
	{
		homeTimelineGroup := v1.Group("/hometimeline")
		{
			htc := &controllers.HomeTimelineController{Metrics: pm, Db: ds.GetDatabase()}
			go homeTimelineGroup.GET("", htc.GetHomeTimeline(ctx))
		}
		userTimelineGroup := v1.Group("/usertimeline")
		{
			utc := &controllers.UserTimelineController{Metrics: pm, Db: ds.GetDatabase()}
			go userTimelineGroup.GET("", utc.GetUserTimeline(ctx))
		}
	}

	return router
}
