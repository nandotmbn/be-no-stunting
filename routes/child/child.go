package routes

import (
	controllers "be-no-stunting-v2/controllers/child"

	"github.com/gin-gonic/gin"
)

func ChildRoute(router *gin.RouterGroup) {
	router.GET("/child/", controllers.ChildHome())
	router.GET("/child/calendar", controllers.ChildCalendar())
	router.GET("/child/measure", controllers.ChildMeasure())
}
