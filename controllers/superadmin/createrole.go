package controllers

import (
	models "be-no-stunting-v2/models"
	"be-no-stunting-v2/views"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

var validateRoles = validator.New()

func CreateRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var roles models.Roles
		defer cancel()

		c.BindJSON(&roles)

		if validationErr := validateRoles.Struct(&roles); validationErr != nil {
			errorMessages := []string{}
			for _, e := range validationErr.(validator.ValidationErrors) {
				errorMessage := fmt.Sprintf("Error on field %s, condition %s = %s", e.Field(), e.ActualTag(), e.Param())
				errorMessages = append(errorMessages, errorMessage)
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": errorMessages})
			return
		}

		filter := bson.M{"$or": []bson.M{ // you can try this in []interface
			{"name": roles.Name},
			{"displayname": roles.DisplayName},
		}}

		count, err := rolesCollection.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		if count >= 1 {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Name has been taken"}})
			return
		}

		newRoles := models.Roles{
			Name:        roles.Name,
			DisplayName: roles.DisplayName,
		}

		result, err := rolesCollection.InsertOne(ctx, newRoles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, views.MasterResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}
