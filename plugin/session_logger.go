package plugin

import (
	"github.com/gin-gonic/gin"
	"time"
	"brief_framework/logger"
)

func LoggerM() gin.HandlerFunc {
	return LoggerWithWriter()
}

func LoggerWithWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		collectData(float64(latency))

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		if method == "POST" {
			logger.GetSessionLogger().Info(" | %3d | %13v | %15s |  %s %-7s  %s %s | return JSON: %s",
				statusCode,
				latency,
				clientIP,
				method,
				path,
				c.GetString("postByte"),
				comment,
				c.GetString("returnString"),
			)
		} else {
			logger.GetSessionLogger().Info(" | %3d | %13v | %15s |  %s %-7s %s | return JSON: %s",
				statusCode,
				latency,
				clientIP,
				method,
				path,
				comment,
				c.GetString("returnString"),
			)
		}
	}
}

func RecoveryM() gin.HandlerFunc {
	return gin.RecoveryWithWriter(logger.Instance())
}
