package utils

import (
	"github.com/gin-gonic/gin"
	"regexp"
)

func ValidateOrderNumber(number string) bool {
	regex := regexp.MustCompile(`^[0-9]+$`)
	return regex.MatchString(number)
}

func GetId(c *gin.Context) (string, bool) {
	userId, exists := c.Get("user_id")
	if !exists {
		return "", exists
	}
	userIdString, _ := userId.(string)
	return userIdString, exists
}
