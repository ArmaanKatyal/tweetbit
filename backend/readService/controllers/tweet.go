package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TweetController struct {
	Metrics *internal.PromMetrics
	DB      *gorm.DB
}

// This flow is triggered when a user clicks on a tweet and wants to see the tweet in detail.
// This will return the tweet and all the replies to the tweet.
func (tc *TweetController) GetTweet(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		_, ok := c.Get("decodedToken")
		if !ok {
			log.Error().Msg("user claims not found in context")
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "user claims not found in context",
			})
			collectMetrics(tc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/tweet",
			}, start)
			return
		}

		tweetId := c.Query("id")
		parsedTweetId, err := strconv.Atoi(tweetId)
		if err != nil {
			log.Error().Msg("tweet id not found in query params")
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "tweet id not found in query params",
			})
			collectMetrics(tc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/tweet",
			}, start)
			return
		}

		if tweetId == "" {
			log.Error().Msg("tweet id not found in query params")
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "tweet id not found in query params",
			})
			collectMetrics(tc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/tweet",
			}, start)
			return
		}

		var tweet models.Tweet
		if err := tc.DB.Table("Tweet").Where(&models.Tweet{Id: uint(parsedTweetId)}).First(&tweet).Error; err != nil {
			log.Error().Msg("tweet not found in database")
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "tweet not found in database",
			})
			collectMetrics(tc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/tweet",
			}, start)
			return
		}

		var comments []models.Tweet_Comments
		if err := tc.DB.Table("Tweet_Comments").Where(&models.Tweet_Comments{Tweet_id: uint(parsedTweetId)}).Find(&comments).Error; err != nil {
			log.Error().Msg("comments not found in database")
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "comments not found in database",
			})
			collectMetrics(tc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/tweet",
			}, start)
			return
		}

		log.Info().Str("module", "tweet.controller").Str("function", "GetTweet").Msg("tweet found in database")
		c.JSON(200, gin.H{
			"tweet":    tweet,
			"comments": comments,
		})
	}
}
