package controllers

import (
	"context"
	"net/http"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/gin-gonic/gin"
)

type UserTimelineController struct {
	Metrics *internal.PromMetrics
}

func (utc *UserTimelineController) GetUserTimeline(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from user timeline"})
	}
}
