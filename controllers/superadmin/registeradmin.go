package controllers

import (
	"be-no-stunting-v2/configs"
	"be-no-stunting-v2/models"
	"be-no-stunting-v2/views"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// collection.
var validateUser = validator.New()

var rolesCollection *mongo.Collection = configs.GetCollection(configs.DB, "roles")
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func RegisteringAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		c.BindJSON(&user)

		if validationErr := validateUser.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		count, err_ := userCollection.CountDocuments(ctx, bson.M{"identifier": user.Identifier})

		if err_ != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Internal server error"}})
			return
		}

		if count >= 1 {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Email or username has been taken"}})
			return
		}

		rolesCount, rolesGetError := rolesCollection.CountDocuments(ctx, bson.M{"_id": user.RolesId})

		if rolesGetError != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Internal server error"}})
			return
		}

		if rolesCount == 0 {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Roles ID not configured"}})
			return
		}

		var choosedRoles views.RolesWithId
		rolesFindOne := rolesCollection.FindOne(ctx, bson.M{"_id": user.RolesId})
		rolesFindOne.Decode(&choosedRoles)

		bytes, errors := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if errors != nil {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Password tidak valid"}})
		}

		newUser := models.User{
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Identifier: user.Identifier,
			Address:    user.Address,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			RolesId:    user.RolesId,
			Password:   string(bytes),
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, bson.M{
			"Status":  http.StatusCreated,
			"Message": "success",
			"Data": bson.M{
				"id":        result.InsertedID,
				"firstname": newUser.FirstName,
				"lastname":  newUser.LastName,
				"Roles": bson.M{
					"id":          choosedRoles.Id,
					"name":        choosedRoles.Name,
					"displayname": choosedRoles.DisplayName,
				},
				"identifier": newUser.Identifier,
			},
		})
	}
}
