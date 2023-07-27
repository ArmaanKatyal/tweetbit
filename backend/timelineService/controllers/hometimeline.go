package controllers

import (
	"context"

	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/internal"
	"github.com/gin-gonic/gin"
)

type HomeTimelineController struct {
	Metrics *internal.PromMetrics
}

func (htc *HomeTimelineController) GetHomeTimeline(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello from home-timeline",
		})
	}
}
