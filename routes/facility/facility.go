package routes

import (
	controllers "be-no-stunting-v2/controllers/facility"

	"github.com/gin-gonic/gin"
)

func FacilityRoute(router *gin.RouterGroup) {
	router.GET("/facility/measure/", controllers.FacilityMeasureFindGet())
	router.POST("/facility/measure/record", controllers.FacilityMeasureRecord())

	router.GET("/facility/monitor", controllers.FacilityMonitorRetrive())
	router.GET("/facility/monitor/:patientId", controllers.FacilityMonitorRetriveByID())
	router.GET("/facility/monitor/:patientId/comment/:postId", controllers.FacilityMonitorRetriveByID())
	// router.POST("/facility/monitor/record", controllers.FacilityMonitorRecord())

	// router.PUT("/master/roles/:roleId", controllers.EditRole())

	// router.GET("/master/roles", controllers.RetriveAllRoles())
	// router.GET("/master/roles/:roleId", controllers.RetriveRole())
	// router.DELETE("/master/roles/", controllers.DeleteAllRoles())
	// router.DELETE("/master/roles/:roleId", controllers.DeleteRole())
}
