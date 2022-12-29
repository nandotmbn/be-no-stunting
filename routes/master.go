package routes

import (
	controllers "be-no-stunting-v2/controllers/master"

	"github.com/gin-gonic/gin"
)

func MasterRoute(router *gin.RouterGroup) {
	router.GET("/master", controllers.GetMasterData())
}
