package main

import (
	"be-no-stunting-v2/configs"
	routes "be-no-stunting-v2/routes"
	routesFacility "be-no-stunting-v2/routes/facility"
	routesPatient "be-no-stunting-v2/routes/patient"
	routesSuperAdmin "be-no-stunting-v2/routes/superadmin"

	// facilitylevel "be-no-stunting/routes/facility/measure"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	configs.ConnectDB()

	// Super Admin Route
	superadmin := router.Group("/superadmin")
	routesSuperAdmin.SuperUserRoute(superadmin)

	// Non Tokenize Route
	v1 := router.Group("/v1")
	routes.AuthRoute(v1)
	routes.MasterRoute(v1)

	// Screen Level Route
	routesFacility.FacilityRoute(v1)
	routesPatient.PatientRoute(v1)

	router.Run("10.252.132.40:6000")
}
