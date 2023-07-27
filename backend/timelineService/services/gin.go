package services

import (
	"context"

	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/controllers"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/middleware"
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
	router.Use(middleware.VerifyToken(pm))

	return router
}
