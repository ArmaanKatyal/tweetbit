package middlewares

import (
	"os"

	"github.com/ArmaanKatyal/tweetbit/backend/authService/helpers"
	"github.com/ArmaanKatyal/tweetbit/backend/authService/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helpers.ExtractAuthToken(c)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized - No Token Provided",
			})
			return
		}

		claims := &models.Claims{}

		parsedJwt, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			HandleJwtErrorResponse(c, err)
		}

		if !parsedJwt.Valid {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized - Invalid Token",
			})
			return
		}

		if claims.Token_type != "access" {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized - Invalid Token",
			})
			return
		}

		c.Set("decodedToken", claims)
	}
}

func VerifyRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helpers.ExtractRefreshToken(c)
		
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized - No Token Provided",
			})
			return
		}

		claims := &models.Claims{}

		parsedJwt, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			HandleJwtErrorResponse(c, err)
		}

		if !parsedJwt.Valid {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized - Invalid Token",
			})
			return
		}

		if claims.Token_type != "refresh" {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized - Invalid Token",
			})
			return
		}

		c.Set("decodedToken", claims)
	}
}

func HandleJwtErrorResponse(c *gin.Context, err error) {
	if err == jwt.ErrSignatureInvalid {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "Unauthorized - Invalid Signature",
		})
		return
	} else if err == jwt.ErrTokenExpired {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "Unauthorized - Token Expired",
		})
		return
	} else if err == jwt.ErrInvalidKey {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "Unauthorized - Invalid Key",
		})
		return
	}
}
