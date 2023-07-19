package services

import (
	"context"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/controllers"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/middlewares"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func NewRouter(ctx context.Context, pm *internal.PromMetrics) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	_, span := otel.Tracer("searchService.services").Start(ctx, "requestRouter")
	defer span.End()

	health := &controllers.HealthController{}
	health.Metrics = pm

	router.GET("/health", health.Status)
	router.GET("/metrics", internal.PrometheusHandler())
	router.Use(middlewares.VerifyToken(pm))

	// v1 := router.Group("/api/v1")
	{
		// searchGroup := v1.Group("/search")
		{
			// span.SetAttributes(attribute.Key("group").String("search"))
			// search := new(controllers.SearchController)
			// search.Metrics = pm
			// searchGroup.GET("/tweet", search.TweetSearch(newCtx))
		}
	}
	return router
}