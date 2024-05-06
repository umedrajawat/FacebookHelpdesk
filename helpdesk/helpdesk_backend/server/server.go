package server

import (

	// "github.com/getsentry/sentry-go"

	configs "helpdesk_backend/config"
	"helpdesk_backend/controllers"
	"helpdesk_backend/db"
	"helpdesk_backend/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var r *gin.Engine

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Init() {
	gin.SetMode(gin.ReleaseMode)

	r = NewRouter()
	// Setting up DBs
	db.SetupRedisDB()
	db.SetupMongoDB()
	db.NewCacheMap()

	r.Use(middlewares.AllowedOriginsMiddleWare())
	r.Use(middlewares.CorsMiddleWare())

	initializeRoutes()

	r.GET("/ws", controllers.HandleWebsocket)

	r.Run(configs.SERVER_PORT)
}
