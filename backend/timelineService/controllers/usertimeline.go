package controllers

import (
	"context"
	"sort"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserTimelineController struct {
	Metrics *internal.PromMetrics
	Db      *gorm.DB
}

func (htc *UserTimelineController) GetUserTimeline(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		userId := c.Query("userId")
		order := c.Query("order")
		if userId == "" {
			log.Error().Msg("userId is required")
			c.JSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "userId is required",
			})
			internal.CollectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/usertimeline",
			}, start)
			return
		}

		// TODO: fetch tweets from postgres
		var tweets []models.Tweet
		err := htc.Db.Table("Tweet").Where("user_id = ?", userId).Find(&tweets).Error
		if err != nil {
			log.Error().Msgf("error while fetching tweets for userId: %s, error: %s", userId, err.Error())
			c.JSON(500, gin.H{
				"error":   "internal_server_error",
				"message": "error while fetching tweets",
			})
			internal.CollectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/usertimeline",
			}, start)
			return
		}

		// sort tweets by created_at
		if order == "desc" {
			sort.SliceStable(tweets, func(i, j int) bool {
				return tweets[i].Created_at.After(*tweets[j].Created_at)
			})
		} else if order == "asc" {
			sort.SliceStable(tweets, func(i, j int) bool {
				return tweets[i].Created_at.Before(*tweets[j].Created_at)
			})
		}

		c.JSON(200, gin.H{
			"tweets": tweets,
		})

	}
}
