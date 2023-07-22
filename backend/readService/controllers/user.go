package controllers

import (
	"context"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/middlewares"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserController struct {
	Metrics *internal.PromMetrics
	DB      *gorm.DB
}

func (uc *UserController) GetUser(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		userClaims, ok := c.Get("decodedToken")
		if !ok {
			log.Error().Msg("user claims not found in context")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "user claims not found in context",
			})
			uc.Metrics.ObserveResponseTime(internal.BadRequest, internal.GET, time.Since(start).Seconds())
			uc.Metrics.IncHttpTransaction(internal.BadRequest, internal.GET)
			return
		}

		// Get user from database
		var user models.User
		uc.DB.Where("email = ?", userClaims.(*middlewares.Claims).Email).Table("User").Unscoped().First(&user)

		// Get user's tweets from database
		var tweets []models.Tweet
		uc.DB.Where("user_id = ?", user.Id).Where("uuid = ?", user.Uuid).Table("Tweet").Unscoped().Find(&tweets)

		log.Info().Str("module", "controller.user").Str("function", "GetUser").Msg("user retrieved successfully")
		c.JSON(200, gin.H{
			"user":   user,
			"tweets": tweets,
		})
		uc.Metrics.ObserveResponseTime(internal.Ok, internal.GET, time.Since(start).Seconds())
		uc.Metrics.IncHttpTransaction(internal.Ok, internal.GET)
	}
}

func (uc *UserController) GetUserLikes(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		userClaims, ok := c.Get("decodedToken")
		if !ok {
			log.Error().Msg("user claims not found in context")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "user claims not found in context",
			})
			uc.Metrics.ObserveResponseTime(internal.BadRequest, internal.GET, time.Since(start).Seconds())
			uc.Metrics.IncHttpTransaction(internal.BadRequest, internal.GET)
			return
		}

		var user models.User
		uc.DB.Where("email = ?", userClaims.(*middlewares.Claims).Email).Table("User").Unscoped().First(&user)

		// Get user liked tweets from database
		var tweets []models.Tweet
		uc.DB.Table("Tweet").Joins(`JOIN "Tweet_Likes" ON "Tweet".id = "Tweet_Likes".tweet_id`).Where(`"Tweet_Likes".user_id = ?`, user.Id).Unscoped().Find(&tweets)

		log.Info().Str("module", "controller.user").Str("function", "GetUserLikes").Msg("user likes retrieved successfully")
		c.JSON(200, gin.H{
			"tweets": tweets,
		})
		uc.Metrics.ObserveResponseTime(internal.Ok, internal.GET, time.Since(start).Seconds())
		uc.Metrics.IncHttpTransaction(internal.Ok, internal.GET)
	}
}

func (uc *UserController) GetUserReplies(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		userClaims, ok := c.Get("decodedToken")
		if !ok {
			log.Error().Msg("user claims not found in context")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "user claims not found in context",
			})
			uc.Metrics.ObserveResponseTime(internal.BadRequest, internal.GET, time.Since(start).Seconds())
			uc.Metrics.IncHttpTransaction(internal.BadRequest, internal.GET)
			return
		}

		var user models.User
		uc.DB.Where("email = ?", userClaims.(*middlewares.Claims).Email).Table("User").Unscoped().First(&user)

		var comments []models.Tweet_Comments
		uc.DB.Raw(`select tc.* from "Tweet_Comments" tc join "Tweet" t on t.id = tc.tweet_id where t.user_id=1;`).Find(&comments)

		var tweets []models.Tweet
		uc.DB.Raw(`select t.* from "Tweet_Comments" tc join "Tweet" t on t.id = tc.tweet_id where t.user_id=1;`).Find(&tweets)

		type reply struct {
			Tweet   models.Tweet
			Comment models.Tweet_Comments
		}

		var replies []reply

		// Merge tweets and comments
		for i := 0; i < len(tweets); i++ {
			replies = append(replies, reply{
				Tweet:   tweets[i],
				Comment: comments[i],
			})
		}

		log.Info().Str("module", "controller.user").Str("function", "GetUserReplies").Msg("user replies retrieved successfully")
		c.JSON(200, gin.H{
			"replies": replies,
		})
		uc.Metrics.ObserveResponseTime(internal.Ok, internal.GET, time.Since(start).Seconds())
		uc.Metrics.IncHttpTransaction(internal.Ok, internal.GET)
	}
}
