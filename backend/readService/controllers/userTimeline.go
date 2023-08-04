package controllers

import (
	"context"
	"encoding/json"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/middlewares"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

type UserTimelineController struct {
	Metrics *internal.PromMetrics
	DB      *gorm.DB
}

func (utc *UserTimelineController) GetUserTimeline(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		userClaims, exists := c.Get("decodedToken")
		if !exists {
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "user claims not found in context",
			})
			collectMetrics(utc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/usertimeline",
			}, start)
			return
		}

		// fetch the user from the database
		var user models.User
		err := utc.DB.Where("email = ?", userClaims.(*middlewares.Claims).Email).Table("User").Unscoped().First(&user).Error
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error retrieving user from database",
			})
			collectMetrics(utc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/usertimeline",
			}, start)
			return
		}

		req, err := http.NewRequest(http.MethodGet, constructUserTimelineUrl(&user), nil)
		req.Header.Set("x-api-key", viper.GetString("timelineservice.apikey"))
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error constructing request to timeline service",
			})
			collectMetrics(utc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/usertimeline",
			}, start)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error making request to timeline service",
			})
			collectMetrics(utc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/usertimeline",
			}, start)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		if resp.StatusCode != 200 {
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "operation_not_allowed",
				"message": "error making request to timeline service",
			})
			collectMetrics(utc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/usertimeline",
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
			collectMetrics(utc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/usertimeline",
			}, start)
			return
		}

		c.JSON(200, gin.H{
			"tweets": response.Tweets,
		})
		collectMetrics(utc.Metrics, &internal.MetricsInput{
			Code:   internal.Ok,
			Method: internal.GET,
			Route:  "/usertimeline",
		}, start)
	}
}

// constructUserTimelineUrl constructs the url to make a request to the timeline service
func constructUserTimelineUrl(m *models.User, order ...string) string {
	if len(order) == 0 {
		order = append(order, "desc")
	}
	return viper.GetString("timelineservice.url") + "usertimeline?userId=" + helpers.UintToString(m.Id) + "&order=" + order[0]
}
