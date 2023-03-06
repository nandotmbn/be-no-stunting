package main

import (
	"be-no-stunting-v2/configs"
	"be-no-stunting-v2/helpers"
	routes "be-no-stunting-v2/routes"
	routesChild "be-no-stunting-v2/routes/child"
	routesFacility "be-no-stunting-v2/routes/facility"
	routesMother "be-no-stunting-v2/routes/mother"
	routesPatient "be-no-stunting-v2/routes/patient"
	routesSuperAdmin "be-no-stunting-v2/routes/superadmin"

	// facilitylevel "be-no-stunting/routes/facility/measure"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	configs.ConnectDB()
	configs.SetupFirebase()

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
	routesChild.ChildRoute(v1)
	routesMother.MotherRoute(v1)

	helpers.RolesSetup()

	router.Run(":8080")
}
