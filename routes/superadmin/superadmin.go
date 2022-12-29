package routes

import (
	controllers "be-no-stunting-v2/controllers/superadmin"

	"github.com/gin-gonic/gin"
)

func SuperUserRoute(router *gin.RouterGroup) {
	router.POST("/admin", controllers.RegisteringAdmin())

	router.POST("/master/roles", controllers.CreateRoles())
	router.PUT("/master/roles/:roleId", controllers.EditRole())
	router.GET("/master/roles", controllers.RetriveAllRoles())
	router.GET("/master/roles/:roleId", controllers.RetriveRole())
	router.DELETE("/master/roles/", controllers.DeleteAllRoles())
	router.DELETE("/master/roles/:roleId", controllers.DeleteRole())
}
