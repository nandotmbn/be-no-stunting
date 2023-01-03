package controllers

import (
	// "be-no-stunting-v2/configs"
	"be-no-stunting-v2/configs"
	// "be-no-stunting-v2/helpers"
	"be-no-stunting-v2/models"
	"fmt"

	"be-no-stunting-v2/views"
	viewsFacility "be-no-stunting-v2/views/facility"
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
		// var record models.Record
		var commentInput viewsFacility.CommentInput
		// var user views.UserNoPassword
		var postIdString = c.Param("postId")
		var patientIdString = c.Param("patientId")
		defer cancel()

		c.BindJSON(&commentInput)

		// var idUser, err = helpers.ValidateToken(helpers.ExtractToken(c))
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError,
		// 		bson.M{
		// 			"Status":  http.StatusInternalServerError,
		// 			"Message": "Internal Server Error",
		// 		},
		// 	)
		// 	return
		// }

		// objId, objIdErr := primitive.ObjectIDFromHex(idUser)
		// if objIdErr != nil {
		// 	c.JSON(http.StatusInternalServerError,
		// 		bson.M{
		// 			"Status":  http.StatusInternalServerError,
		// 			"Message": "Internal Server Error",
		// 		},
		// 	)
		// 	return
		// }

		if validationErr := validate.Struct(&commentInput); validationErr != nil {
			errorMessages := []string{}
			for _, e := range validationErr.(validator.ValidationErrors) {
				errorMessage := fmt.Sprintf("Error on field %s, condition %s = %s", e.Field(), e.ActualTag(), e.Param())
				errorMessages = append(errorMessages, errorMessage)
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": errorMessages})
			return
		}

		var postId, errPostId = primitive.ObjectIDFromHex(postIdString)
		if errPostId != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "PostId is not valid ObjectID"})
			return
		}

		var patientId, errPatientId = primitive.ObjectIDFromHex(patientIdString)
		if errPatientId != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "PatientId is not valid ObjectID"})
			return
		}

		countMonitor, errCountMonitor := monitorCollection.CountDocuments(ctx, bson.M{"_id": postId, "patientid": patientId})
		countRecord, errCountRecord := recordCollection.CountDocuments(ctx, bson.M{"_id": postId, "patientid": patientId})
		commentCount, errCommentCount := commentCollection.CountDocuments(ctx, bson.M{"postid": postId})
		if errCountMonitor != nil || errCountRecord != nil || (countMonitor < 1 && countRecord < 1) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "There is no post by given ID"})
			return
		}
		if errCommentCount != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Internal Server Error"})
			return
		}
		if commentCount > 0 {
			update := bson.M{
				"content": commentInput.Content,
			}

			updateComment, err := commentCollection.UpdateOne(ctx, bson.M{"postid": postId}, bson.M{"$set": update})
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			c.JSON(http.StatusCreated, updateComment)
			return
		}

		newComment := models.Comment{
			PostId:    postId,
			Content:   commentInput.Content,
			CreatedAt: time.Now(),
		}

		update := bson.M{
			"ischecked": true,
		}
		if countMonitor > 0 {
			_, err := commentCollection.UpdateOne(ctx, bson.M{"_id": postId}, bson.M{"$set": update})
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		} else if countRecord > 0 {
			_, err := recordCollection.UpdateOne(ctx, bson.M{"_id": postId}, bson.M{"$set": update})
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		result, err := commentCollection.InsertOne(ctx, newComment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, views.MasterResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}
