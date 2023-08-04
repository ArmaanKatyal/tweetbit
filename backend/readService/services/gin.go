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
		homeTimelineGroup := v1.Group("/hometimeline")
		{
			htc := controllers.HomeTimelineController{Metrics: pm, DB: ds.GetDatabase()}
			go homeTimelineGroup.GET("", htc.GetHomeTimeline(ctx))
		}
		userTimelineGroup := v1.Group("/usertimeline")
		{
			utc := controllers.UserTimelineController{Metrics: pm, DB: ds.GetDatabase()}
			go userTimelineGroup.GET("", utc.GetUserTimeline(ctx))
		}
		userGroup := v1.Group("/user")
		{
			uc := controllers.UserController{Metrics: pm, DB: ds.GetDatabase()}
			go userGroup.GET("", uc.GetUser(ctx))
			go userGroup.GET("/likes", uc.GetUserLikes(ctx))
			go userGroup.GET("/replies", uc.GetUserReplies(ctx))
			go userGroup.GET("/followers", uc.GetUserFollowers(ctx))
			go userGroup.GET("/following", uc.GetUserFollowing(ctx))
		}
		tweetGroup := v1.Group("/tweet")
		{
			tc := controllers.TweetController{Metrics: pm, DB: ds.GetDatabase()}
			go tweetGroup.GET("", tc.GetTweet(ctx))
		}
	}
	return router
}
