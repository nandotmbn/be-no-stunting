package controllers

import (
	"be-no-stunting-v2/configs"
	"be-no-stunting-v2/helpers"
	"be-no-stunting-v2/models"
	"be-no-stunting-v2/views"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var fcmtokenCollection *mongo.Collection = configs.GetCollection(configs.DB, "fcmtoken")

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var inputLogin views.InputLogin
		defer cancel()

		c.BindJSON(&inputLogin)

		var resultLogin models.User
		var finalUserView views.UserNoPassword
		result := userCollection.FindOne(ctx, bson.M{"identifier": inputLogin.Identifier})
		result.Decode(&resultLogin)
		result.Decode(&finalUserView)
		err := bcrypt.CompareHashAndPassword([]byte(resultLogin.Password), []byte(inputLogin.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, bson.M{
				"Status":  http.StatusBadRequest,
				"Message": "Bad request",
				"Data":    "Identifier or Password is not valid",
			})
			return
		}

		var choosedRoles views.RolesWithId
		rolesFindOne := rolesCollection.FindOne(ctx, bson.M{"_id": resultLogin.RolesId})
		rolesFindOne.Decode(&choosedRoles)

		jwtResult, jwtError := helpers.GenerateJWT(resultLogin.Id.String(), choosedRoles.Name)
		if jwtError != nil {

			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		count, err_ := fcmtokenCollection.CountDocuments(ctx, bson.M{"userid": resultLogin.Id, "fcmtoken": inputLogin.FCMToken})

		if err_ != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Internal server error"}})
			return
		}

		if count >= 1 {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "This device with this account has been login"}})
			return
		}

		newLoginFCM := models.FCMToken{
			UserId:   resultLogin.Id,
			FCMToken: inputLogin.FCMToken,
		}

		resultFCM, err := fcmtokenCollection.InsertOne(ctx, newLoginFCM)
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var FCMToken models.FCMToken
		resultFCMFinal := fcmtokenCollection.FindOne(ctx, bson.M{"_id": resultFCM.InsertedID})
		resultFCMFinal.Decode(&FCMToken)

		c.JSON(http.StatusOK,
			bson.M{
				"Status":   http.StatusOK,
				"Message":  "Success",
				"Data":     finalUserView,
				"Token":    jwtResult,
				"FCMToken": FCMToken,
			},
		)
	}
}
