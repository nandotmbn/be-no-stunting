package controllers

import (
	"be-no-stunting-v2/configs"
	"be-no-stunting-v2/helpers"
	"fmt"

	views "be-no-stunting-v2/views"
	viewsFacility "be-no-stunting-v2/views/facility"
	"context"

	// "fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var rolesCollection *mongo.Collection = configs.GetCollection(configs.DB, "roles")

// Retrive single user using by its ID
func FacilityMeasureFindGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user []viewsFacility.FacilityMeasureFindGet
		defer cancel()
		var paramName = c.Query("name")

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

		var childRole views.RolesWithId

		__err := rolesCollection.FindOne(ctx, bson.M{"name": "Child"}).Decode(&childRole)
		if __err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		results, err := userCollection.Find(ctx, bson.M{
			"parentid": idUser,
			"rolesid":  childRole.Id,
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
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer results.Close(ctx)

		for results.Next(ctx) {
			var singleRoles viewsFacility.FacilityMeasureFindGet
			if err = results.Decode(&singleRoles); err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			user = append(user, singleRoles)
		}

		c.JSON(http.StatusOK, bson.M{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    user,
		})
	}
}
