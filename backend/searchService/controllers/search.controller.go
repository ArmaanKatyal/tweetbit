package controllers

import (
	"context"
	"net/http"
	"time"

	elastic "github.com/ArmaanKatyal/tweetbit/backend/searchService/elasticsearch"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/internal"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SearchController struct {
	Metrics *internal.PromMetrics
}

func (sc SearchController) TweetSearch(ctx context.Context, client *elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		newCtx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("searchService.controllers").Start(ctx, "SearchController.Search")
		defer span.End()
		query := c.Query("q")

		span.SetAttributes(attribute.Key("query").String(query))
		span.SetAttributes(attribute.Key("path").String(c.Request.URL.Path))
		span.SetAttributes(attribute.Key("method").String(c.Request.Method))
		span.SetAttributes(attribute.Key("host").String(c.Request.Host))
		span.SetAttributes(attribute.Key("user_agent").String(c.Request.UserAgent()))
		span.SetAttributes(attribute.Key("client_ip").String(c.ClientIP()))

		if query == "" {
			span.SetAttributes(attribute.Key("error").Bool(true))
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
			return
		}

		// base64 decode the query
		decodedQuery, err := helpers.Base64Decode(query)
		if err != nil {
			sc.Metrics.ObserveResponseTime(internal.InternalServerError, internal.GET, time.Since(start).Seconds())
			sc.Metrics.IncHttpTransaction(internal.InternalServerError, internal.GET)
			span.SetAttributes(attribute.Key("error").Bool(true))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		esTweet := elastic.NewElasticTweet(client)
		tweets, err := esTweet.TweetSearch(newCtx, &decodedQuery)
		if err != nil {
			sc.Metrics.ObserveResponseTime(internal.InternalServerError, internal.GET, time.Since(start).Seconds())
			sc.Metrics.IncHttpTransaction("500", "GET")
			span.SetAttributes(attribute.Key("error").Bool(true))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		sc.Metrics.ObserveResponseTime(internal.Ok, internal.GET, time.Since(start).Seconds())
		sc.Metrics.IncHttpTransaction(internal.Ok, internal.GET)
		c.JSON(http.StatusOK, gin.H{"tweets": tweets})
	}
}
