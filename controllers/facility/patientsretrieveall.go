package controllers

import (
	"be-no-stunting-v2/helpers"

	views "be-no-stunting-v2/views"
	"context"
	"fmt"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Retrive single user using by its ID
func FacilityPatientRetrive() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var paramName = c.Query("name")
		var paramRole = c.Query("role")

		var userIdString, err = helpers.ValidateToken(helpers.ExtractToken(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		userId, err := primitive.ObjectIDFromHex(userIdString)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}
		var patientTypeId primitive.ObjectID

		var resultsFinal *mongo.Cursor
		if len(paramRole) > 0 {
			roleId, err2 := primitive.ObjectIDFromHex(paramRole)
			if err2 != nil {
				c.JSON(http.StatusInternalServerError,
					bson.M{
						"Status":  http.StatusInternalServerError,
						"Message": "Internal Server Error",
					},
				)
				return
			}
			patientTypeId = roleId
			results, err := userCollection.Find(ctx, bson.M{
				"parentid": userId,
				"rolesid":  patientTypeId,
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

			resultsFinal = results
		} else {
			results, err := userCollection.Find(ctx, bson.M{
				"parentid": userId,
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

			resultsFinal = results
		}

		var resultMonitor []views.UserNoPassword
		for resultsFinal.Next(ctx) {
			var singleRoles views.UserNoPassword
			if err = resultsFinal.Decode(&singleRoles); err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			resultMonitor = append(resultMonitor, singleRoles)
		}

		defer resultsFinal.Close(ctx)

		c.JSON(200, bson.M{
			"Status":  200,
			"Message": "Success",
			"Data":    resultMonitor,
		})

	}
}
