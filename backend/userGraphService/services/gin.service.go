package services

import (
	"net/http"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/userGraphService/internal"
	"github.com/gin-gonic/gin"
)

func NewRouter(pm *internal.PromMetrics) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		start := time.Now()
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
		pm.IncHttpTransaction(internal.Ok, internal.GET)
		pm.ObserveResponseTime(internal.Ok, internal.GET, time.Since(start).Seconds())
	})

	router.GET("/metrics", internal.PrometheusHandler())

	return router
}
