package routes

import (
	controllers "be-no-stunting-v2/controllers/patients"

	"github.com/gin-gonic/gin"
)

func PatientRoute(router *gin.RouterGroup) {
	router.POST("/patient/monitor/record", controllers.FacilityMonitorRecord())
	// router.PUT("/master/roles/:roleId", controllers.EditRole())

	// router.GET("/master/roles/:roleId", controllers.RetriveRole())
	// router.DELETE("/master/roles/", controllers.DeleteAllRoles())
	// router.DELETE("/master/roles/:roleId", controllers.DeleteRole())
}
