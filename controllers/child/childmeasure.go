package controllers

import (
	// "be-no-stunting-v2/configs"

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

// Retrive single user using by its ID
func ChildMeasure() gin.HandlerFunc {
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

		matchRecordAgg := bson.D{
			{
				Key: "$match", Value: bson.M{"patientid": objId},
			},
		}
		groupRecordStage := bson.D{
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

		cursorRecord, errCursorRecord := recordCollection.Aggregate(ctx, mongo.Pipeline{matchRecordAgg, groupStage, groupRecordStage})
		if errCursorRecord != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorRecord.Error()}})
			return
		}

		var resultRecord []viewsFacility.FacilityRecordHome
		if errCursorRecord = cursorRecord.All(ctx, &resultRecord); errCursorRecord != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorRecord.Error()}})
			return
		}

		c.JSON(http.StatusOK, bson.M{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data": bson.M{
				"child":   user,
				"measure": resultRecord,
			},
		})
	}
}
