package controllers

import (
	"context"

	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/internal"
	"github.com/gin-gonic/gin"
)

type UserTimelineController struct {
	Metrics *internal.PromMetrics
}

func (htc *UserTimelineController) GetUserTimeline(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello from user-timeline",
		})
	}
}
