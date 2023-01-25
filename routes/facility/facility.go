package routes

import (
	controllers "be-no-stunting-v2/controllers/facility"

	"github.com/gin-gonic/gin"
)

func FacilityRoute(router *gin.RouterGroup) {
	router.GET("/facility/measure/", controllers.FacilityMeasureFindGet())
	router.POST("/facility/measure/record", controllers.FacilityMeasureRecord())
	router.GET("/facility/measure/patient/:patientId", controllers.FacilityMonitorGetUserData())

	router.GET("/facility/monitor/calendar", controllers.FacilityMonitorRetrive())
	router.GET("/facility/monitor/calendar/:patientId", controllers.FacilityMonitorRetriveByID())
	router.GET("/facility/monitor/calendar/:patientId/comment/:postId", controllers.FacilityMonitorCommentGet())
	router.POST("/facility/monitor/calendar/:patientId/comment/:postId", controllers.FacilityMonitorCommentPost())
	router.GET("/facility/monitor/calendar/:patientId/comment/:postId/check", controllers.FacilityMonitorCheck())

	router.GET("/facility/monitor/record", controllers.FacilityMeasureRetrive())
	router.GET("/facility/monitor/record/:patientId/comment/:postId/check", controllers.FacilityMeasureCheck())

	// router.PUT("/master/roles/:roleId", controllers.EditRole())

	// router.GET("/master/roles", controllers.RetriveAllRoles())
	// router.GET("/master/roles/:roleId", controllers.RetriveRole())
	// router.DELETE("/master/roles/", controllers.DeleteAllRoles())
	// router.DELETE("/master/roles/:roleId", controllers.DeleteRole())
}
