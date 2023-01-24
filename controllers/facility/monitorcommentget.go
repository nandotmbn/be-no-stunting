package controllers

import (
	// "be-no-stunting-v2/configs"

	// "be-no-stunting-v2/helpers"

	views "be-no-stunting-v2/views/facility"
	"context"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Retrive single user using by its ID
func FacilityMonitorCommentGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var postIdString = c.Param("postId")
		var patientIdString = c.Param("patientId")
		defer cancel()

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

		var finalResult views.CommentInput
		result := commentCollection.FindOne(ctx, bson.M{"postid": postId, "userid": patientId})
		result.Decode(&finalResult)

		c.JSON(http.StatusOK, bson.M{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    finalResult,
		})
	}
}
