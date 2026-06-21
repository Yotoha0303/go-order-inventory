package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	RequestIDHeader = "X-Request-ID"
	RequestKeyID    = "request_id"
)

type RequestIDContextKey = struct{}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)

		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(RequestKeyID, requestID)

		ctx := context.WithValue(
			c.Request.Context(),
			RequestIDContextKey{},
			requestID,
		)

		c.Request = c.Request.WithContext(ctx)

		c.Header(RequestKeyID, requestID)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	value, exists := c.Get(RequestKeyID)
	if !exists {
		return ""
	}

	requestID, _ := value.(string)
	return requestID
}

func RequestIDFromContext(c context.Context) string {
	requestID, _ := c.Value(RequestIDContextKey{}).(string)
	return requestID
}
