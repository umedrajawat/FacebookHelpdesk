package middlewares

import (
	"helpdesk_backend/logger"

	"github.com/gin-gonic/gin"
)

func UpdateMissingHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.ZapLogger.Info(c.Request.Header.Get("Content-Type"))
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.Request.Header.Set("Content-Type", "application/json")
		}

		c.Next()
	}
}
