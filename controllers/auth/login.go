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

		jwtResult, jwtError := helpers.GenerateJWT(resultLogin.Id, choosedRoles.Name)
		if jwtError != nil {

			c.JSON(http.StatusInternalServerError,
				bson.M{
					"Status":  http.StatusInternalServerError,
					"Message": "Internal Server Error",
				},
			)
			return
		}

		c.JSON(http.StatusOK,
			bson.M{
				"Status":  http.StatusOK,
				"Message": "Success",
				"Data":    finalUserView,
				"Token":   jwtResult,
			},
		)
	}
}
