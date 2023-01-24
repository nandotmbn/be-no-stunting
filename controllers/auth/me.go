package controllers

import (
	"be-no-stunting-v2/helpers"
	"be-no-stunting-v2/views"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var paramId = c.Param("id")
		defer cancel()

		userId, err := primitive.ObjectIDFromHex(paramId)
		if err != nil {
			c.JSON(http.StatusBadRequest, bson.M{
				"Status":  http.StatusBadRequest,
				"Message": "Bad request",
				"Data":    "Id that you sent is invalid",
			})
			return
		}

		var idParentString, err_ = helpers.ValidateToken(helpers.ExtractToken(c))
		if err_ != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		var idParent, err__ = primitive.ObjectIDFromHex(idParentString)
		if err__ != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		var finalResult views.UserNoPassword
		result := userCollection.FindOne(ctx, bson.M{"_id": userId, "parentid": idParent})
		result.Decode(&finalResult)

		c.JSON(http.StatusOK,
			bson.M{
				"Status":  http.StatusOK,
				"Message": "Success",
				"Data":    finalResult,
			},
		)
	}
}
