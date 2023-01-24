package controllers

import (
	// "be-no-stunting-v2/helpers"

	views "be-no-stunting-v2/views"
	viewsFacility "be-no-stunting-v2/views/facility"
	"context"
	"fmt"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Retrive single user using by its ID
func FacilityMonitorCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var commentInput viewsFacility.CommentInput
		var postIdString = c.Param("postId")
		var patientIdString = c.Param("patientId")
		defer cancel()

		c.BindJSON(&commentInput)

		if validationErr := validate.Struct(&commentInput); validationErr != nil {
			errorMessages := []string{}
			for _, e := range validationErr.(validator.ValidationErrors) {
				errorMessage := fmt.Sprintf("Error on field %s, condition %s = %s", e.Field(), e.ActualTag(), e.Param())
				errorMessages = append(errorMessages, errorMessage)
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": errorMessages})
			return
		}

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

		countMonitor, errCountMonitor := monitorCollection.CountDocuments(ctx, bson.M{"_id": postId, "patientid": patientId})
		if errCountMonitor != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "There is no post by given ID"})
			return
		}

		var monitor viewsFacility.FacilityMonitorGet
		monitorResult := monitorCollection.FindOne(ctx, bson.M{"_id": postId, "patientid": patientId})
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
			_, err := monitorCollection.UpdateOne(ctx, bson.M{"_id": postId}, bson.M{"$set": update})
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, bson.M{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    monitor,
		})
	}
}
