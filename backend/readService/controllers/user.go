package controllers

import (
	"context"

	"github.com/ArmaanKatyal/tweetbit/backend/readService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/middlewares"
	"github.com/ArmaanKatyal/tweetbit/backend/readService/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Metrics     *internal.PromMetrics
	Dataservice *internal.DatabaseService
}

func (uc *UserController) GetUser(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {

		userClaims, ok := c.Get("decodedToken")
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{
				"message": "Internal server error",
			})
			return
		}

		var user models.User
		uc.Dataservice.Db.Where("email = ?", userClaims.(*middlewares.Claims).Email).Table("User").Unscoped().First(&user)

		c.JSON(200, gin.H{
			"user": user,
		})
	}
}
