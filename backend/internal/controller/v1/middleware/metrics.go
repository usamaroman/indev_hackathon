package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/metrics"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := getStatusCodeLabel(c.Writer.Status())
		path := c.FullPath()

		metrics.IncHttpRequestsTotal(c.Request.Method, path, status)
		metrics.SetHttpRequestDuration(c.Request.Method, path, duration)
		metrics.SetHttpResponseSize(c.Request.Method, path, float64(c.Writer.Size()))
	}
}

func getStatusCodeLabel(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	case status >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}
