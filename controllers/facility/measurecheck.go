package controllers

import (
	// "be-no-stunting-v2/helpers"

	views "be-no-stunting-v2/views"
	viewsFacility "be-no-stunting-v2/views/facility"
	"context"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Retrive single user using by its ID
func FacilityMeasureCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var postIdString = c.Param("postId")
		var patientIdString = c.Param("patientId")
		defer cancel()

		var postId, errPostId = primitive.ObjectIDFromHex(postIdString)
		if errPostId != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "PostId is not valid ObjectID"})
			return
		}

		var patientId, errPatientId = primitive.ObjectIDFromHex(patientIdString)
		if errPatientId != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "PatientId is not valid ObjectID"})
			return
		}

		countMonitor, errCountMonitor := recordCollection.CountDocuments(ctx, bson.M{"_id": postId, "patientid": patientId})
		if errCountMonitor != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "There is no post by given ID"})
			return
		}

		var monitor viewsFacility.FacilityMonitorGet
		monitorResult := recordCollection.FindOne(ctx, bson.M{"_id": postId, "patientid": patientId})
		monitorResult.Decode(&monitor)

		var update primitive.M = bson.M{
			"ischecked": false,
		}

		if !monitor.IsChecked {
			update = bson.M{
				"ischecked": true,
			}
		}

		if countMonitor > 0 {
			_, err := recordCollection.UpdateOne(ctx, bson.M{"_id": postId}, bson.M{"$set": update})
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		var monitorDone viewsFacility.FacilityMonitorGet
		monitorResults := recordCollection.FindOne(ctx, bson.M{"_id": postId, "patientid": patientId})
		monitorResults.Decode(&monitorDone)

		c.JSON(http.StatusOK, bson.M{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    monitorDone,
		})
	}
}
