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
		err = uc.DB.Raw(`select tc.* from "Tweet_Comments" tc join "Tweet" t on t.id = tc.tweet_id where t.user_id= ?;`, user.Id).Find(&comments).Error
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
		err = uc.DB.Raw(`select t.* from "Tweet_Comments" tc join "Tweet" t on t.id = tc.tweet_id where t.user_id= ?;`, user.Id).Find(&tweets).Error
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

func (uc *UserController) GetUserFollowers(ctx context.Context) gin.HandlerFunc {
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
				Route:  "/user/followers",
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
				Route:  "/user/followers",
			}, start)
			return
		}

		var followers []models.User
		err = uc.DB.Raw(`select u.* from "User" u join "User_Followers" uf on u.id = uf.follower_id where uf.user_id = ?;`, user.Id).Find(&followers).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user followers from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user followers from database",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/user/followers",
			}, start)
			return
		}

		log.Info().Str("module", "controller.user").Str("function", "GetUserFollowers").Msg("user followers retrieved successfully")
		c.JSON(200, gin.H{
			"followers": followers,
		})
		collectMetrics(uc.Metrics, &internal.MetricsInput{
			Code:   internal.Ok,
			Method: internal.GET,
			Route:  "/user/followers",
		}, start)
	}
}

func (uc *UserController) GetUserFollowing(ctx context.Context) gin.HandlerFunc {
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
				Route:  "/user/following",
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
				Route:  "/user/following",
			}, start)
			return
		}

		var following []models.User
		err = uc.DB.Raw(`select u.* from "User" u join "User_Followers" uf on u.id = uf.user_id where uf.follower_id = ?;`, user.Id).Find(&following).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user following from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user following from database",
			})
			collectMetrics(uc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/user/following",
			}, start)
			return
		}

		log.Info().Str("module", "controller.user").Str("function", "GetUserFollowing").Msg("user following retrieved successfully")
		c.JSON(200, gin.H{
			"following": following,
		})
		collectMetrics(uc.Metrics, &internal.MetricsInput{
			Code:   internal.Ok,
			Method: internal.GET,
			Route:  "/user/following",
		}, start)
	}
}

func collectMetrics(pm *internal.PromMetrics, metrics *internal.MetricsInput, t time.Time) {
	pm.ObserveResponseTime(metrics.Code, metrics.Method, metrics.Route, time.Since(t).Seconds())
	pm.IncHttpTransaction(metrics.Code, metrics.Method, metrics.Route)
}
