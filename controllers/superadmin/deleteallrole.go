package controllers

import (
	"be-no-stunting-v2/views"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteAllRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := rolesCollection.DeleteMany(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				views.MasterResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Role with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			views.MasterResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Role successfully deleted!"}},
		)
	}
}
