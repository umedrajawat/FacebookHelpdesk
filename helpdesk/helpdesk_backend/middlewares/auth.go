package middlewares

import (
	"fmt"
	configs "helpdesk_backend/config"
	"helpdesk_backend/db"
	"helpdesk_backend/logger"
	"helpdesk_backend/utilities"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware returns a Gin middleware handler function
func AuthMiddleware(ctx *gin.Context) {
	user := ctx.GetHeader("x-user")
	authToken := strings.Split(ctx.GetHeader("x-authToken"), " ")[1]
	fetchedToken := db.GetRedisKey(db.RedisClient, user)
	logger.ZapLogger.Info("x-user ", user)
	logger.ZapLogger.Info("x-authToken ", authToken)
	logger.ZapLogger.Info("fetched token ", fetchedToken)
	if (fetchedToken != "") && fetchedToken == authToken {
		logger.ZapLogger.Info("user successfully authenticated")
		ctx.Next()
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
			"status":    "error",
			"error_msg": "unauthorized request",
		})
		return
	}
}

func AllowedOriginsMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		headerOrigin := c.Request.Header.Get("Origin")
		hostParam := c.Query("host")

		headerOrigin = strings.Replace(headerOrigin, "https://", "", -1)
		headerOrigin = strings.Replace(headerOrigin, "http://", "", -1)
		headerOrigin = strings.Replace(headerOrigin, "/", "", -1)

		ok := utilities.Contains(configs.ALLOWED_HOSTS, headerOrigin)

		fmt.Println("ok", ok, headerOrigin)
		if !ok && hostParam != "" {
			ok = utilities.Contains(configs.ALLOWED_HOSTS, hostParam)
		}
		logger.ZapLogger.Info("IS OK", ok)
		if !ok {
			c.String(400, "Goofy Origin")
			c.Abort()
		}

		c.Next()
	}
}
