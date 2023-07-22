package services

import (
	"context"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/controllers"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/middlewares"
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
	router.Use(middlewares.VerifyToken(pm))

	v1 := router.Group("/api/v1")
	{
		homeTimelineGroup := v1.Group("/home_timeline")
		{
			htc := controllers.HomeTimelineController{Metrics: pm}
			homeTimelineGroup.GET("", htc.GetHomeTimeline(ctx))
		}
		userTimelineGroup := v1.Group("/user_timeline")
		{
			htc := controllers.UserTimelineController{Metrics: pm}
			userTimelineGroup.GET("", htc.GetUserTimeline(ctx))
		}
		userGroup := v1.Group("/user")
		{
			uc := controllers.UserController{Metrics: pm, DB: ds.GetDatabase()}
			userGroup.GET("", uc.GetUser(ctx))
			userGroup.GET("/likes", uc.GetUserLikes(ctx))
			userGroup.GET("/replies", uc.GetUserReplies(ctx))
			userGroup.GET("/followers", uc.GetUserFollowers(ctx))
			userGroup.GET("/following", uc.GetUserFollowing(ctx))
		}
	}
	return router
}
