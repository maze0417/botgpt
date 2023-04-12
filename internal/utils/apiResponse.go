package utils

import (
	"botgpt/internal/models"
	"github.com/gin-gonic/gin"
)

func getRequestId(c *gin.Context) string {
	value, exist := c.Get("requestId")

	if !exist {
		return ""
	}
	return value.(string)
}

func SendResponse(code int, result models.ResponseResult, c *gin.Context) {
	result.RequestId = getRequestId(c)

	c.JSON(code, result)
}
func SendPageResponse(code int, result models.ResponsePageResult, c *gin.Context) {
	result.RequestId = getRequestId(c)
	c.JSON(code, result)
}
