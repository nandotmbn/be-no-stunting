package controllers

import (
	// "be-no-stunting-v2/configs"

	"be-no-stunting-v2/configs"
	"be-no-stunting-v2/helpers"

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

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var monitorCollection *mongo.Collection = configs.GetCollection(configs.DB, "monitor")

// Retrive single user using by its ID
func MotherHome() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// var monitor models.Monitor
		var user views.UserNoPassword
		var facility views.UserNoPassword
		var paramTime = c.Query("datetime")
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
		userCollection.FindOne(ctx, bson.M{"_id": user.ParentId}).Decode(&facility)

		matchMonitorAgg := bson.D{
			{
				Key: "$match", Value: bson.M{"patientid": objId},
			},
		}
		var dateAgg primitive.D

		if len(paramTime) < 1 {
			dateAgg = matchMonitorAgg
		} else {
			parsedTime, err := time.Parse(time.RFC3339, paramTime)
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}

			dateAgg = bson.D{
				{
					Key: "$match", Value: bson.M{"createdat": bson.M{
						"$gte": parsedTime,
					}},
				},
			}
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

		cursorMonitor, errCursorMonitor := monitorCollection.Aggregate(ctx, mongo.Pipeline{matchMonitorAgg, dateAgg, groupMonitorStage})
		if errCursorMonitor != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorMonitor.Error()}})
			return
		}

		var resultMonitor []viewsFacility.FacilityMonitorFindByIdView
		if errCursorMonitor = cursorMonitor.All(ctx, &resultMonitor); errCursorMonitor != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorMonitor.Error()}})
			return
		}

		c.JSON(http.StatusOK, bson.M{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data": bson.M{
				"facility": facility,
				"mother":   user,
				"monitor":  resultMonitor,
			},
		})
	}
}
