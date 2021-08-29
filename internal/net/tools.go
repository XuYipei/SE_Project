package net

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context) string {
	result := c.Request.Header.Get("Authorization")
	for i := 0; i < len(result); i++ {
		if result[i] == ' ' {
			return strings.SplitN(result, " ", 2)[1]
		}
	}
	return result
}
