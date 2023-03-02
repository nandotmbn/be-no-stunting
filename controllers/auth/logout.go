package controllers

import (
	"be-no-stunting-v2/views"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var inputLogout views.InputLogout
		defer cancel()

		c.BindJSON(&inputLogout)

		objId, _ := primitive.ObjectIDFromHex(inputLogout.UserId)

		count, err_ := fcmtokenCollection.CountDocuments(ctx, bson.M{"userid": objId, "fcmtoken": inputLogout.FCMToken})

		if err_ != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Internal server error"}})
			return
		}

		if count == 0 {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "This device with this account has not been login"}})
			return
		}

		result, err := fcmtokenCollection.DeleteOne(ctx, bson.M{"userid": objId, "fcmtoken": inputLogout.FCMToken})
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		fmt.Println(result)

		c.JSON(http.StatusOK,
			bson.M{
				"Status":  http.StatusOK,
				"Message": "Success",
				"Data":    "Successfully Logout",
			},
		)
	}
}
