package controllers

import (
	"net/http"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/gin-gonic/gin"
)

type HealthController struct {
	Metrics *internal.PromMetrics
}

func (hc *HealthController) Status(c *gin.Context) {
	start := time.Now()
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
	hc.Metrics.ObserveResponseTime(internal.Success, internal.GET, "/health", time.Since(start).Seconds())
	hc.Metrics.IncHttpTransaction(internal.Success, internal.GET, "/health")
}
