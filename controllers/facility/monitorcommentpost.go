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

var commentCollection *mongo.Collection = configs.GetCollection(configs.DB, "comments")

// Retrive single user using by its ID
func FacilityMonitorCommentPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var record models.Record
		var user views.UserNoPassword
		defer cancel()

		c.BindJSON(&record)
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

		objId, _ := primitive.ObjectIDFromHex(idUser)
		if validationErr := validate.Struct(&record); validationErr != nil {
			errorMessages := []string{}
			for _, e := range validationErr.(validator.ValidationErrors) {
				errorMessage := fmt.Sprintf("Error on field %s, condition %s = %s", e.Field(), e.ActualTag(), e.Param())
				errorMessages = append(errorMessages, errorMessage)
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": errorMessages})
			return
		}

		userCollection.FindOne(ctx, bson.M{"_id": record.PatientId}).Decode(&user)

		if user.ParentId != objId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You cannot record unregistered patient"})
			return
		}

		newRecord := models.Record{
			Height:     record.Height,
			Weight:     record.Weight,
			HeartRate:  record.HeartRate,
			Temp:       record.Temp,
			PatientId:  record.PatientId,
			FacilityId: objId,
			CreatedAt:  time.Now(),
		}

		result, err := recordCollection.InsertOne(ctx, newRecord)
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, views.MasterResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}