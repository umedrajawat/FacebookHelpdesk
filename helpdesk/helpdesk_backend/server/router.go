package server

import (
	"helpdesk_backend/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// router.Use(middlewares.AuthMiddleware())

	// v1 := router.Group("v1")
	// {
	// 	userGroup := v1.Group("user")
	// 	{
	// 		user := new(controllers.UserController)
	// 		userGroup.GET("/:id", user.Retrieve)
	// 	}
	// }
	return router

}

func initializeRoutes() {
	health := new(controllers.HealthController)
	r.GET(Health, health.Status)
	r.GET(Webhook, controllers.Webhook)
	r.POST(Webhook, controllers.Webhook)
	r.POST(CreateUserr, controllers.CreateUser)
	r.GET(GetMsgs, controllers.GetAllMessagesController)
	r.GET(GetUser, controllers.GetUser)
	r.POST(SendMsg, controllers.SendMesaage)
	r.POST(LoginUrl, controllers.LoginUser) // resp done

	// authRouter := r.Group("", middlewares.AuthMiddleware)
	// {
	// 	// authRouter.POST(LoginUrl, controllers.LoginUser) // resp done
	// }

}
