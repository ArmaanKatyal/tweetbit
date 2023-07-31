package controllers

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type HomeTimelineController struct {
	Metrics *internal.PromMetrics
	Db      *gorm.DB
}

func (htc *HomeTimelineController) GetHomeTimeline(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		userId := c.Query("userId")
		if userId == "" {
			log.Error().Str("module", "controllers.hometimeline").Str("function", "GetHomeTimeline").Msg("userId is required")
			c.JSON(400, gin.H{
				"error":   "operation_not_allowed",
				"message": "userId is required",
			})
			internal.CollectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.BadRequest,
				Method: internal.GET,
				Route:  "/hometimeline",
			}, start)
			return
		}

		// fetch hometimeline from redis
		tweets, err := homeTimelineFromRedis(ctx, userId)
		if err != nil {
			log.Error().Str("module", "controllers.hometimeline").Str("function", "GetHomeTimeline").Msgf("error while fetching tweets for userId: %s, error: %s", userId, err.Error())
			c.JSON(500, gin.H{
				"error":   "internal_server_error",
				"message": "error while fetching tweets",
			})
			internal.CollectMetrics(htc.Metrics, &internal.MetricsInput{
				Code:   internal.InternalServerError,
				Method: internal.GET,
				Route:  "/hometimeline",
			}, start)
			return
		}

		// sort tweets by likes_count in descending order
		sort.SliceStable(tweets, func(i, j int) bool {
			return tweets[i].Likes_count > tweets[j].Likes_count
		})

		c.JSON(200, gin.H{
			"tweets": tweets,
		})
		internal.CollectMetrics(htc.Metrics, &internal.MetricsInput{
			Code:   internal.Ok,
			Method: internal.GET,
			Route:  "/hometimeline",
		}, start)
	}
}

func homeTimelineFromRedis(ctx context.Context, userId string) ([]models.Tweet, error) {
	rdb := internal.NewRedisServer()
	tweetClient := rdb.GetTweetClient()
	defer rdb.Close()
	tweetsList, err := tweetClient.LRange(ctx, userId, 0, -1).Result()
	if err != nil {
		log.Error().Str("module", "controllers.hometimeline").Str("function", "homeTimelineFromRedis").Msgf("error while fetching tweets for userId: %s, error: %s", userId, err.Error())
		return nil, err
	}

	var tweets []models.Tweet
	// parse string to models.Tweet
	for _, tweet := range tweetsList {
		t, err := parseTweet(tweet)
		if err != nil {
			log.Error().Str("module", "controllers.hometimeline").Str("function", "homeTimelineFromRedis").Msgf("error while unmarshalling tweet: %s, error: %s", tweet, err.Error())
			return nil, err
		}
		tweets = append(tweets, t)
	}
	return tweets, nil
}

func parseTweet(tweet string) (models.Tweet, error) {
	var t models.Tweet_Redis
	err := json.Unmarshal([]byte(tweet), &t)
	if err != nil {
		log.Error().Str("module", "controllers.hometimeline").Str("function", "parseTweet").Msgf("error while unmarshalling tweet: %s, error: %s", tweet, err.Error())
		return models.Tweet{}, err
	}

	parsedTweet := models.Tweet{
		Id:             helpers.ConvertStringToUint(t.Id),
		Uuid:           t.Uuid,
		User_id:        helpers.ConvertStringToUint(t.User_id),
		Content:        t.Content,
		Created_at:     t.Created_at,
		Likes_count:    helpers.ConvertStringToUint(t.Likes_count),
		Retweets_count: helpers.ConvertStringToUint(t.Retweets_count),
	}
	return parsedTweet, nil
}
