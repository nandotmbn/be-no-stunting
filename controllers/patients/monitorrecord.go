package controllers

import (
	// "be-no-stunting-v2/configs"
	"be-no-stunting-v2/configs"
	"be-no-stunting-v2/helpers"
	"be-no-stunting-v2/models"
	"fmt"

	"be-no-stunting-v2/views"
	"context"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var monitorCollection *mongo.Collection = configs.GetCollection(configs.DB, "monitor")
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

// Retrive single user using by its ID
func FacilityMonitorRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var monitor models.Monitor
		var user models.User
		defer cancel()

		c.BindJSON(&monitor)
		var idUser, err = helpers.ValidateToken(helpers.ExtractToken(c))

		if err != nil {
			panic(err)
		}

		if validationErr := validate.Struct(&monitor); validationErr != nil {
			errorMessages := []string{}
			for _, e := range validationErr.(validator.ValidationErrors) {
				errorMessage := fmt.Sprintf("Error on field %s, condition %s = %s", e.Field(), e.ActualTag(), e.Param())
				errorMessages = append(errorMessages, errorMessage)
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
			return
		}

		objId, _ := primitive.ObjectIDFromHex(idUser)

		err__ := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err__ != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		newMonitor := models.Monitor{
			IsChecked:     false,
			PatientTypeId: user.RolesId,
			PatientId:     objId,
			Content:       monitor.Content,
			FacilityId:    user.ParentId,
			CreatedAt:     time.Now(),
		}

		result, err := monitorCollection.InsertOne(ctx, newMonitor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		update := bson.M{
			"updatedat": time.Now(),
		}
		userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		c.JSON(http.StatusCreated, views.MasterResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}
