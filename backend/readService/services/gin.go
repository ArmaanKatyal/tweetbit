package services

import (
	"context"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/controllers"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/middlewares"
	"github.com/gin-gonic/gin"
)

func NewRouter(ctx context.Context, pm *internal.PromMetrics) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := &controllers.HealthController{}
	health.Metrics = pm

	router.GET("/health", health.Status)
	router.GET("/metrics", internal.PrometheusHandler())
	router.Use(middlewares.VerifyToken(pm))

	v1 := router.Group("/api/v1")
	{
		homeTimelineGroup := v1.Group("/home_timeline")
		{
			htc := controllers.HomeTimelineController{}
			htc.Metrics = pm
			homeTimelineGroup.GET("", htc.GetHomeTimeline(ctx))
		}
		userTimelineGroup := v1.Group("/user_timeline")
		{
			htc := controllers.UserTimelineController{}
			htc.Metrics = pm
			userTimelineGroup.GET("", htc.GetUserTimeline(ctx))
		}
	}
	return router
}
