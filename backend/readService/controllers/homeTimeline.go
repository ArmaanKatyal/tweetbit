package controllers

import (
	"context"
	"net/http"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/gin-gonic/gin"
)

type HomeTimelineController struct {
	Metrics *internal.PromMetrics
}

func (htc *HomeTimelineController) GetHomeTimeline(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from home timeline"})
	}
}
