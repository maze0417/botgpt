package middleware

import (
	"botgpt/internal/config"
	res "botgpt/internal/enum"
	"botgpt/internal/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"strings"
)

func ExceptionMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {

			detailError := fmt.Sprintf("Exception occurs: InternalServerError %v", err)

			if config.IsNotProduction() {
				detailError = fmt.Sprintf("Exception occurs: %v %s", err, parseStack())
			}

			log.Error(detailError)

			c.AbortWithStatusJSON(http.StatusInternalServerError,
				response.Failure(detailError, res.FALIURE))

		}
	}()
	c.Next()
}

func parseStack() string {
	var file string
	var line int

	// Skip the first 3 frames, as they are not relevant to the error location
	const skip = 3
	pc := make([]uintptr, 32)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	var builder strings.Builder

	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		file = frame.File
		line = frame.Line
		builder.WriteString(fmt.Sprintf("=> %s %d ", file, line))

	}

	return builder.String()
}
