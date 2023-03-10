package routes

import (
	controllers "be-no-stunting-v2/controllers/auth"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.RouterGroup) {
	router.POST("/auth/register", controllers.Register())
	router.POST("/auth/login", controllers.Login())
	router.POST("/auth/logout", controllers.Logout())
	router.GET("/auth/me/:id", controllers.Me())
	router.PUT("/auth/me/:id", controllers.Update())
}
