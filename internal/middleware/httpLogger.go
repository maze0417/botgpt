package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
	"time"
)

var skipLogPath = []string{
	"/api/v1/captcha/getCaptcha",
	"/api/v1/merchant/uploadSendMutipleMessage",
	"/api/v1/merchant/uploadSendMutipleMessageOnlyNumber",
	"/api/v1/merchant/uploadSendMutipleMessageOnlyNumber",
	"/conversation",
}

func HttpLoggerMiddleware(c *gin.Context) {

	//response.SetRequestContext(c)
	requestId := c.Request.Header.Get("X-Request-Id")

	if requestId == "" {
		requestId = uuid.New().String()
		c.Request.Header.Set("X-Request-Id", requestId)
	}

	c.Set("requestId", requestId)

	startTime := time.Now()

	// 在 HTTP Header 中設置 requestId
	blw := &bodyLogWriter{
		body:           bytes.NewBufferString(""),
		ResponseWriter: c.Writer,
	}
	c.Writer = blw

	requestBody, _ := c.GetRawData()
	jsonBody, err := strconv.Unquote(strings.Replace(strconv.Quote(string(requestBody)), `\\u`, `\u`, -1))
	if err != nil {
		// handle error
		jsonBody = string(requestBody)
	}

	// 将请求 body 再写回到 gin.Context 中，以便后续的处理器可以使用
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	// 處理 Request
	defer func() {

		endTime := time.Now()

		// 記錄日誌
		statusCode := blw.Status()
		duration := endTime.Sub(startTime)

		responseBody := blw.body.String()

		if contains(skipLogPath, c.Request.RequestURI) {
			responseBody = "skip log ..."
		}

		logString := fmt.Sprintf("[HTTP] %v %v %v %v %v", requestId, c.Request.Method, c.Request.URL.Path, statusCode, duration)
		log.WithFields(log.Fields{
			"type":           "http",
			"requestId":      requestId,
			"requestHeaders": getRequestHeaders(c),
			"requestUrl":     c.Request.RequestURI,
			"requestBody":    jsonBody,
			"statusCode":     statusCode,
			"duration":       duration.Seconds(),
			"responseBody":   responseBody,
		}).Info(logString)
	}()
	c.Header("Content-Type", "application/json")
	c.Next()

}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

// getRequestHeaders 從請求中獲取標頭並返回一個易讀的字串
func getRequestHeaders(c *gin.Context) string {
	var headers strings.Builder
	for key, values := range c.Request.Header {
		for _, value := range values {
			headers.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}
	return headers.String()
}
