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
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "user claims not found in context",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/user",
			}, start)
			return
		}

		// Get user from database
		var user models.User
		err := uc.DB.Where("email = ?", userClaims.(*middlewares.Claims).Email).Table("User").Unscoped().First(&user).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user from database",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/user",
			}, start)
			return
		}

		// Get user's tweets from database
		var tweets []models.Tweet
		err = uc.DB.Where("user_id = ?", user.Id).Where("uuid = ?", user.Uuid).Table("Tweet").Unscoped().Find(&tweets).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user tweets from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user tweets from database",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/user",
			}, start)
			return
		}

		log.Info().Str("module", "controller.user").Str("function", "GetUser").Msg("user retrieved successfully")
		c.JSON(200, gin.H{
			"user":   user,
			"tweets": tweets,
		})
		collectMetrics(uc.Metrics, &internal.MetricsInput{
			Code:   internal.Ok,
			Method: internal.GET,
			Route:  "/user",
		}, start)
	}
}

func (uc *UserController) GetUserLikes(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		userClaims, ok := c.Get("decodedToken")
		if !ok {
			log.Error().Msg("user claims not found in context")
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "user claims not found in context",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/user/likes",
			}, start)
			return
		}

		var user models.User
		err := uc.DB.Where("email = ?", userClaims.(*middlewares.Claims).Email).Table("User").Unscoped().First(&user).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user from database",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/user/likes",
			}, start)
			return
		}

		// Get user liked tweets from database
		var tweets []models.Tweet
		err = uc.DB.Table("Tweet").Joins(`JOIN "Tweet_Likes" ON "Tweet".id = "Tweet_Likes".tweet_id`).Where(`"Tweet_Likes".user_id = ?`, user.Id).Unscoped().Find(&tweets).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user liked tweets from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user liked tweets from database",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/user/likes",
			}, start)
			return
		}

		log.Info().Str("module", "controller.user").Str("function", "GetUserLikes").Msg("user likes retrieved successfully")
		c.JSON(200, gin.H{
			"tweets": tweets,
		})
		collectMetrics(uc.Metrics, &internal.MetricsInput{
			Code:   internal.Ok,
			Method: internal.GET,
			Route:  "/user/likes",
		}, start)
	}
}

func (uc *UserController) GetUserReplies(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		userClaims, ok := c.Get("decodedToken")
		if !ok {
			log.Error().Msg("user claims not found in context")
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "user claims not found in context",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/user/replies",
			}, start)
			return
		}

		var user models.User
		err := uc.DB.Where("email = ?", userClaims.(*middlewares.Claims).Email).Table("User").Unscoped().First(&user).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user from database",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/user/replies",
			}, start)
			return
		}

		var comments []models.Tweet_Comments
		err = uc.DB.Raw(`select tc.* from "Tweet_Comments" tc join "Tweet" t on t.id = tc.tweet_id where t.user_id=1;`).Find(&comments).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user comments from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user comments from database",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/user/replies",
			}, start)
			return
		}

		var tweets []models.Tweet
		err = uc.DB.Raw(`select t.* from "Tweet_Comments" tc join "Tweet" t on t.id = tc.tweet_id where t.user_id=1;`).Find(&tweets).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user tweets from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user tweets from database",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/user/replies",
			}, start)
			return
		}

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

		collectMetrics(uc.Metrics, &internal.MetricsInput{
			Code:   internal.Ok,
			Method: internal.GET,
			Route:  "/user/replies",
		}, start)
	}
}

func collectMetrics(pm *internal.PromMetrics, metrics *internal.MetricsInput, t time.Time) {
	pm.ObserveResponseTime(metrics.Code, metrics.Method, time.Since(t).Seconds())
	pm.IncHttpTransaction(metrics.Code, metrics.Method)
}
