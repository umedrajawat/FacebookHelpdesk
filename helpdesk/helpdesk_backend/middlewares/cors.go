package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// var CSRFMiddelWare = gincsrf.Middleware(gincsrf.Options{
// 	Secret: configs.SECRET,
// 	ErrorFunc: func(c *gin.Context) {

// 		c.String(400, "CSRF token mismatch")
// 		c.Abort()
// 	},
// 	IgnoreMethods: []string{"GET", "HEAD", "OPTIONS"},
// })

// CorsMiddleWare : Handle CorsMiddleWare
func CorsMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "authorization, x-csrftoken, content-type, x-app-uuid")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "authorization, x-csrftoken, content-type, x-app-uuid")

		// SECURITY HEADER
		c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
		c.Writer.Header().Set("X-XSS-Protection", "1")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func RedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "/") {
			c.Redirect(http.StatusOK, strings.TrimSuffix(c.Request.URL.Path, "/"))
			c.Abort()
			return
		}

		c.Next()
	}
}
