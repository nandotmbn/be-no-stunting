package routes

import (
	controllers "be-no-stunting-v2/controllers/mother"

	"github.com/gin-gonic/gin"
)

func MotherRoute(router *gin.RouterGroup) {
	router.GET("/mother/", controllers.MotherHome())
}
