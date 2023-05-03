package middleware

import (
	"botgpt/internal/clients/chatgpt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//goland:noinspection GoUnhandledErrorResult
func CheckHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(chatgpt.AuthorizationHeader) == "" && c.Request.URL.Path != "/auth/login" {
			c.AbortWithStatusJSON(http.StatusForbidden, chatgpt.ReturnMessage("Missing accessToken."))
			return
		}

		c.Header("Content-Type", "application/json")
		c.Next()
	}
}
