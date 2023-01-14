package controllers

import (
	views "be-no-stunting-v2/views"
	"context"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Retrive single user using by its ID
func FacilityMonitorGetUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var patientIdString = c.Param("patientId")
		var patientId, patientIdErr = primitive.ObjectIDFromHex(patientIdString)
		if patientIdErr != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "PatientId is not valid ObjectID",
				},
			)
			return
		}

		var resultLogin views.UserNoPassword
		result := userCollection.FindOne(ctx, bson.M{"_id": patientId})
		result.Decode(&resultLogin)

		c.JSON(http.StatusOK, bson.M{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    resultLogin,
		})
	}
}
