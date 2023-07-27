package middleware

import (
	"net/http"

	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/internal"
	"github.com/ArmaanKatyal/tweetbit/backend/timelineService/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id         string `json:"id"`
	Email      string `json:"email"`
	Token_type string `json:"type"`
	jwt.RegisteredClaims
}

func VerifyToken(pm *internal.PromMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helpers.ExtractAuthToken(c)
		if token == "" {
			pm.IncHttpTransaction(internal.BadRequest, internal.GET, internal.VerifyToken)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized - No token provided",
			})
			return
		}

		claims := &Claims{}
		parsedJwt, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(utils.GetEnvValue("JWT_SECRET")), nil
		})

		if err != nil {
			pm.IncHttpTransaction(internal.BadRequest, internal.GET, internal.VerifyToken)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		if !parsedJwt.Valid {
			pm.IncHttpTransaction(internal.BadRequest, internal.GET, internal.VerifyToken)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized - Invalid token",
			})
			return
		}

		if claims.Token_type != "access" {
			pm.IncHttpTransaction(internal.BadRequest, internal.GET, internal.VerifyToken)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized - Invalid token type",
			})
			return
		}

		c.Set("decodedToken", claims)
	}
}
