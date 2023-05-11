package middlewares

import (
	"github.com/ArmaanKatyal/tweetbit/backend/authService/helpers"
	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helpers.ExtractAuthToken(c)
		c.Set("decoded", token)
	}
}
