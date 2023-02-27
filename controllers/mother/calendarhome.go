package controllers

import (
	// "be-no-stunting-v2/configs"

	"be-no-stunting-v2/helpers"
	"fmt"

	"be-no-stunting-v2/views"
	viewsFacility "be-no-stunting-v2/views/facility"
	"context"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Retrive single user using by its ID
func MotherCalendar() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// var monitor models.Monitor
		var user views.UserNoPassword
		defer cancel()

		var idUser, err = helpers.ValidateToken(helpers.ExtractToken(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		objId, _ := primitive.ObjectIDFromHex(idUser)

		userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)

		matchMonitorAgg := bson.D{
			{
				Key: "$match", Value: bson.M{"patientid": objId},
			},
		}

		groupMonitorStage := bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "as", Value: "patient"},
				{Key: "localField", Value: "patientid"},
				{Key: "foreignField", Value: "_id"},
			},
			},
		}

		groupStage := bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "comments"},
				{Key: "as", Value: "comment"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "postid"},
			},
			},
		}

		cursorMonitor, errCursorMonitor := monitorCollection.Aggregate(ctx, mongo.Pipeline{matchMonitorAgg, groupStage, groupMonitorStage})
		if errCursorMonitor != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorMonitor.Error()}})
			return
		}

		var resultMonitor []viewsFacility.FacilityMonitorFindByIdView
		if errCursorMonitor = cursorMonitor.All(ctx, &resultMonitor); errCursorMonitor != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorMonitor.Error()}})
			return
		}

		fmt.Println(resultMonitor)

		c.JSON(http.StatusOK, bson.M{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data": bson.M{
				"mother":  user,
				"monitor": resultMonitor,
			},
		})
	}
}
