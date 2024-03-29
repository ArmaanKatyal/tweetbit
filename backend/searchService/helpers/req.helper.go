package helpers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractAuthToken(c *gin.Context) string {

	if value, err := c.Cookie("access_token"); err == nil {
		return value
	}

	return strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
}

func ExtractRefreshToken(c *gin.Context) string {
	return strings.Replace(c.GetHeader("Refresh"), "Bearer ", "", 1)
}
