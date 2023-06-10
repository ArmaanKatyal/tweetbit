package services

import (
	"context"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/controllers"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/middlewares"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewRouter(ctx context.Context, es *elasticsearch.Client) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	newCtx, span := otel.Tracer("searchService.services").Start(ctx, "requestRouter")
	defer span.End()

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	router.Use(middlewares.VerifyToken())

	v1 := router.Group("/api/v1")
	{
		searchGroup := v1.Group("/search")
		{
			span.SetAttributes(attribute.Key("group").String("search"))
			search := new(controllers.SearchController)
			searchGroup.GET("/tweet", search.TweetSearch(newCtx, es))
		}
	}
	return router
}
