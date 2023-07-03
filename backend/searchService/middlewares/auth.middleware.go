package middlewares

import (
	"net/http"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(pm *internal.PromMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helpers.ExtractAuthToken(c)
		if token == "" {
			pm.IncHttpTransaction(internal.BadRequest, internal.GET)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized - No token provided",
			})
			return
		}

		claims := &models.Claims{}
		parsedJwt, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(helpers.GetEnvValue("JWT_SECRET")), nil
		})

		if err != nil {
			pm.IncHttpTransaction(internal.BadRequest, internal.GET)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		if !parsedJwt.Valid {
			pm.IncHttpTransaction(internal.BadRequest, internal.GET)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized - Invalid token",
			})
			return
		}

		if claims.Token_type != "access" {
			pm.IncHttpTransaction(internal.BadRequest, internal.GET)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized - Invalid token type",
			})
			return
		}

		c.Set("decodedToken", claims)
	}
}
