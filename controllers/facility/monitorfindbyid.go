package controllers

import (
	// "be-no-stunting-v2/helpers"

	"be-no-stunting-v2/helpers"
	views "be-no-stunting-v2/views"
	viewsFacility "be-no-stunting-v2/views/facility"
	"context"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Retrive single user using by its ID
func FacilityMonitorRetriveByID() gin.HandlerFunc {
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

		var idFacilityString, err_ = helpers.ValidateToken(helpers.ExtractToken(c))
		if err_ != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Token is invalid",
				},
			)
			return
		}
		var idFacility, err = primitive.ObjectIDFromHex(idFacilityString)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		matchMonitorAgg := bson.D{
			{
				Key: "$match", Value: bson.M{"patientid": patientId, "facilityid": idFacility},
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

		matchRecordAgg := bson.D{
			{
				Key: "$match", Value: bson.M{"patientid": patientId, "facilityid": idFacility},
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

		cursorMonitor, errCursorMonitor := monitorCollection.Aggregate(ctx, mongo.Pipeline{matchMonitorAgg, groupMonitorStage})
		if errCursorMonitor != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorMonitor.Error()}})
			return
		}

		var resultMonitor []viewsFacility.FacilityMonitorFindByIdView
		if errCursorMonitor = cursorMonitor.All(ctx, &resultMonitor); errCursorMonitor != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorMonitor.Error()}})
			return
		}

		cursorRecord, errCursorRecord := recordCollection.Aggregate(ctx, mongo.Pipeline{matchRecordAgg, groupRecordStage})
		if errCursorRecord != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorRecord.Error()}})
			return
		}

		var resultRecord []viewsFacility.FacilityRecordFindById
		if errCursorRecord = cursorRecord.All(ctx, &resultRecord); errCursorRecord != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": errCursorRecord.Error()}})
			return
		}

		var finalResult views.UserNoPassword
		result := userCollection.FindOne(ctx, bson.M{"_id": patientId, "parentid": idFacility})
		result.Decode(&finalResult)

		c.JSON(http.StatusOK, bson.M{
			"Status":  http.StatusOK,
			"Message": "success",
			"User":    finalResult,
			"Record":  resultRecord,
			"Monitor": resultMonitor,
		})
	}
}
