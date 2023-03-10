package controllers

import (
	"be-no-stunting-v2/helpers"
	"be-no-stunting-v2/views"
	"context"

	// "fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user views.UserUpdate
		c.BindJSON(&user)

		var paramId = c.Param("id")

		userId, err := primitive.ObjectIDFromHex(paramId)
		if err != nil {
			c.JSON(http.StatusBadRequest, bson.M{
				"Status":  http.StatusBadRequest,
				"Message": "Bad request",
				"Data":    "Id that you sent is invalid",
			})
			return
		}

		if validationErr := validateUser.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		var parentIdString, err_ = helpers.ValidateToken(helpers.ExtractToken(c))
		if err_ != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		parentId, err__ := primitive.ObjectIDFromHex(parentIdString)
		if err__ != nil {
			c.JSON(http.StatusBadRequest, bson.M{
				"Status":  http.StatusBadRequest,
				"Message": "Bad request",
				"Data":    "Id that you sent is invalid",
			})
			return
		}
		count, err___ := userCollection.CountDocuments(ctx, bson.M{"_id": userId, "parentid": parentId})
		if err___ != nil {
			c.JSON(404, bson.M{
				"status":  404,
				"message": "Not Found",
			})

			return
		}
		if count == 0 {
			c.JSON(404, bson.M{
				"status":  404,
				"message": "Not Found",
			})

			return
		}

		identifierSameCount, identifierSameErr := userCollection.CountDocuments(ctx, bson.M{"identifier": user.Identifier})
		if identifierSameErr != nil {
			c.JSON(404, bson.M{
				"status":  404,
				"message": "Not Found",
			})

			return
		}
		update := bson.M{
			"firstname":  user.FirstName,
			"lastname":   user.LastName,
			"identifier": user.Identifier,
			"address":    user.Address,
			"rolesid":    user.RolesId,
			"updatedat":  time.Now(),
			"bornat":     user.BornAt,
			"ismale":     user.IsMale,
		}

		if identifierSameCount == 1 {
			var currentIdentifier views.UserOnlyId
			identifierSame := userCollection.FindOne(ctx, bson.M{"identifier": user.Identifier, "parentid": parentId})
			identifierSame.Decode(&currentIdentifier)

			if currentIdentifier.Id == "" {
				c.JSON(http.StatusForbidden,
					bson.M{
						"Status":  http.StatusForbidden,
						"Message": "Try to use other identifier",
					},
				)
				return
			}

			currentIdentifierObjectId, currentIdentifierObjectIdErr := primitive.ObjectIDFromHex(currentIdentifier.Id)
			if currentIdentifierObjectIdErr != nil {
				c.JSON(http.StatusInternalServerError,
					bson.M{
						"Status":  http.StatusInternalServerError,
						"Message": "Internal Server Error",
					},
				)
				return
			}

			if userId != currentIdentifierObjectId {
				c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Identifier has been taken"}})
				return
			}
			result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userId, "parentid": parentId}, bson.M{"$set": update})
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}

			c.JSON(http.StatusAccepted, bson.M{
				"status":  http.StatusAccepted,
				"message": result,
				"data": bson.M{
					"firstName":  user.FirstName,
					"lastName":   user.LastName,
					"rolesId":    user.RolesId,
					"identifier": user.Identifier,
					"address":    user.Address,
				},
			})

			return
		}

		result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userId, "parentid": parentId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusAccepted, bson.M{
			"Status":  http.StatusAccepted,
			"Message": result,
			"Data": bson.M{
				"firstName":  user.FirstName,
				"lastName":   user.LastName,
				"rolesId":    user.RolesId,
				"identifier": user.Identifier,
				"address":    user.Address,
			},
		})
	}
}
