package service

import (
	"net/http"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/internal"
	"github.com/gin-gonic/gin"
)

func NewRouter(pm *internal.PromMetrics) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	router.GET("/metrics", internal.PrometheusHandler())

	return router
}
