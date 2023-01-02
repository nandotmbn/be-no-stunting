package controllers

import (
	"be-no-stunting-v2/configs"
	"be-no-stunting-v2/helpers"
	"strconv"

	views "be-no-stunting-v2/views"
	viewsFacility "be-no-stunting-v2/views/facility"
	"context"
	"fmt"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var monitorCollection *mongo.Collection = configs.GetCollection(configs.DB, "monitor")

// Retrive single user using by its ID
func FacilityMonitorRetrive() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var paramName = c.Query("name")
		var paramChecked = c.Query("checked")
		var paramType = c.Query("type")
		var paramTime = c.Query("datetime")
		var idUserString, err = helpers.ValidateToken(helpers.ExtractToken(c))
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}
		idUser, err := primitive.ObjectIDFromHex(idUserString)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}
		var user []primitive.ObjectID
		ojLorem, err := primitive.ObjectIDFromHex("000000000000000000000000")
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}
		user = append(user, ojLorem)

		results, err := userCollection.Find(ctx, bson.M{
			"$or": []bson.M{
				{
					"firstname": bson.D{
						{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%s.*$", paramName), Options: "si"}},
					},
				},
				{
					"lastname": bson.D{
						{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%s.*$", paramName), Options: "si"}},
					},
				},
			},
		},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer results.Close(ctx)

		for results.Next(ctx) {
			var singleUser views.UserOnlyId
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}

			objIdSingleUser, err := primitive.ObjectIDFromHex(singleUser.Id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}

			user = append(user, objIdSingleUser)
		}

		idUserAgg := bson.D{
			{
				Key: "$match", Value: bson.M{"patientid": bson.M{"$in": user}},
			},
		}
		if paramChecked == "" {
			paramChecked = "false"
		}

		boolValue, err := strconv.ParseBool(paramChecked)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		checkedAgg := bson.D{
			{
				Key: "$match", Value: bson.M{"ischecked": boolValue},
			},
		}

		// create group stage
		groupStage := bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "as", Value: "patient"},
				{Key: "localField", Value: "patientid"},
				{Key: "foreignField", Value: "_id"},
			},
			},
		}

		var dateAgg primitive.D

		if len(paramTime) < 1 {
			dateAgg = checkedAgg
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
						"$lt":  parsedTime.AddDate(0, 0, 1),
					}},
				},
			}
		}

		facilityAgg := bson.D{
			{
				Key: "$match", Value: bson.M{"facilityid": idUser},
			},
		}

		var patientTypeId = ojLorem
		if len(paramType) > 0 {
			patientTypeIdConv, err := primitive.ObjectIDFromHex(paramType)
			if err != nil {
				c.JSON(http.StatusInternalServerError,
					bson.M{
						"Status":  http.StatusInternalServerError,
						"Message": "Internal Server Error",
					},
				)
				return
			}
			patientTypeId = patientTypeIdConv
		}

		patientTypeAgg := bson.D{
			{
				Key: "$match", Value: bson.M{"patienttypeid": patientTypeId},
			},
		}

		// pass the pipeline to the Aggregate() method
		cursor, err := monitorCollection.Aggregate(ctx, mongo.Pipeline{patientTypeAgg, facilityAgg, dateAgg, idUserAgg, groupStage, checkedAgg})
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// display the results
		var resultMonitor []viewsFacility.FacilityMonitorGet
		if err = cursor.All(ctx, &resultMonitor); err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, resultMonitor)
	}
}
