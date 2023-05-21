package controllers

import (
	"os"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/authService/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	// for now just return a auth and refresh token
	// later on we will use the user service to get the user
	// and then return the user and the tokens
	var json LoginBody
	if c.BindJSON(&json) == nil {
		if json.Email == "test@test.com" && json.Password == "test" {
			expirationTime := time.Now().Add(time.Hour)
			claims := &models.Claims{
				Email:      json.Email,
				Token_type: "access",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(expirationTime),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenKey := []byte(os.Getenv("JWT_SECRET"))
			tokenString, err := token.SignedString(tokenKey)

			if err != nil {
				c.JSON(500, gin.H{
					"message": "Internal Server Error",
				})
				return
			}

			c.JSON(200, gin.H{
				"auth_token": tokenString,
			})
		} else {
			c.JSON(401, gin.H{
				"message": "Invalid Credentials",
			})
		}
	}
}

func Logout(c *gin.Context) {
}

func Refresh(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "you got here somehow",
	})
}

func Register(c *gin.Context) {
}
