package controllers

import (
	"net/http"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/internal"
	"github.com/gin-gonic/gin"
)

type HealthController struct {
	Metrics *internal.PromMetrics
}

func (h HealthController) Status(c *gin.Context) {
	start := time.Now()
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
	h.Metrics.ObserveResponseTime(internal.Success, internal.GET, time.Since(start).Seconds())
	h.Metrics.IncHttpTransaction(internal.Success, internal.GET)
}
