package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/middlewares"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type HomeTimelineController struct {
	Metrics *internal.PromMetrics
	DB      *gorm.DB
}

func (htc *HomeTimelineController) GetHomeTimeline(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// make a get request to the timeline-service to get the home timeline
		// return the response

		userClaims, exists := c.Get("decodedToken")
		if !exists {
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "user claims not found in context",
			})
			collectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/hometimeline",
			}, start)
			return
		}

		var user models.User
		err := htc.DB.Where("email = ?", userClaims.(*middlewares.Claims).Email).Table("User").Unscoped().First(&user).Error
		if err != nil {
			log.Error().Err(err).Msg("error retrieving user from database")
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user from database",
			})
			collectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/hometimeline",
			}, start)
			return
		}

		req, err := http.NewRequest(http.MethodGet, constructHomeTimelineUrl(&user), nil)
		req.Header.Set("x-api-key", viper.GetString("timelineservice.apikey"))
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "internal_server_error",
				"message": "unable to make request to timeline service",
			})
			collectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/hometimeline",
			}, start)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "internal_server_error",
				"message": "unable to make request to timeline service",
			})
			collectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/hometimeline",
			}, start)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error().Str("module", "hometimeline.controllers").Str("function", "GetHomeTimeline").Err(err).Msg("error closing response body")
			}
		}(resp.Body)

		if resp.StatusCode != 200 {
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "internal_server_error",
				"message": fmt.Sprintf("unable to make request to timeline service, status code: %d", resp.StatusCode),
			})
			collectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/hometimeline",
			}, start)
			return
		}

		var response struct {
			Tweets []models.Tweet `json:"tweets"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "internal_server_error",
				"message": "unable to decode response from timeline service",
			})
			collectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/hometimeline",
			}, start)
			return
		}

		c.JSON(200, gin.H{
			"tweets": response.Tweets,
		})
		collectMetrics(htc.Metrics, &internal.MetricsInput{
			Code:   internal.Ok,
			Method: internal.GET,
			Route:  "/hometimeline",
		}, start)
	}
}

// constructHomeTimelineUrl constructs the url to make a request to the timeline service
func constructHomeTimelineUrl(user *models.User) string {
	url := viper.GetString("timelineservice.url") + "hometimeline?userId=" + helpers.UintToString(user.Id)
	return url
}
