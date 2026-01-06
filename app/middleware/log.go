package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap" // Assuming the logger is zap.SugaredLogger
)

// CustomResponseWriter wraps gin.ResponseWriter to capture the response body.
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the body before writing to the original ResponseWriter.
func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware creates a Gin middleware for development logging.
// It logs request and response details in a human-readable format.
func LoggerMiddleware(logger *zap.SugaredLogger, skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()

		// --- Log Request ---
		var requestBody string
		contentType := c.GetHeader("Content-Type")

		// Safely read request body and restore it for subsequent handlers
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body
		}

		if strings.Contains(contentType, "multipart/form-data") {
			requestBody = "[Multipart Form Data]"
		} else {
			// Try to pretty-print JSON, otherwise just print as string
			var prettyJSON bytes.Buffer
			if len(bodyBytes) > 0 && json.Indent(&prettyJSON, bodyBytes, "", "  ") == nil {
				requestBody = prettyJSON.String()
			} else {
				requestBody = string(bodyBytes)
			}
		}

		logger.Infof("--> %s %s\n    IP: %s\n    Body: %s", method, path, ip, requestBody)

		// --- Log Response (deferred) ---
		// Use CustomResponseWriter to capture response body
		crw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = crw

		defer func() {
			latency := time.Since(startTime)
			statusCode := crw.Status()

			var responseBody string
			// Try to pretty-print JSON, otherwise just print as string
			var prettyJSON bytes.Buffer
			if len(crw.body.Bytes()) > 0 && json.Indent(&prettyJSON, crw.body.Bytes(), "", "  ") == nil {
				responseBody = prettyJSON.String()
			} else {
				responseBody = crw.body.String()
			}

			// Format latency for better readability
			latencyStr := latency.String()
			if latency > time.Minute { // Format long latencies
				latencyStr = fmt.Sprintf("%.2fm", latency.Minutes())
			} else if latency > time.Second {
				latencyStr = fmt.Sprintf("%.2fs", latency.Seconds())
			} else if latency > time.Millisecond {
				latencyStr = fmt.Sprintf("%.2fms", float64(latency.Nanoseconds())/float64(time.Millisecond))
			} else {
				latencyStr = fmt.Sprintf("%.2fus", float64(latency.Nanoseconds())/float64(time.Microsecond))
			}

			logger.Infof("<-- %s %s %d %s\n    Response: %s", method, path, statusCode, latencyStr, responseBody)
		}()

		c.Next()
	}
}
