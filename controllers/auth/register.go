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
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validateUser = validator.New()

var rolesCollection *mongo.Collection = configs.GetCollection(configs.DB, "roles")

func Register() gin.HandlerFunc {
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
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Identifier has been taken"}})
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

		var userAccount views.UserNoPassword
		accountUserId, _ := primitive.ObjectIDFromHex(idUser)
		userAccountResult := userCollection.FindOne(ctx, bson.M{"_id": accountUserId})
		userAccountResult.Decode(&userAccount)

		var userRolesAccount views.RolesWithId
		userRolesAccountResult := rolesCollection.FindOne(ctx, bson.M{"_id": userAccount.RolesId})
		userRolesAccountResult.Decode(&userRolesAccount)

		var targetRolesAccount views.RolesWithId
		TargetRolesAccountResult := rolesCollection.FindOne(ctx, bson.M{"_id": user.RolesId})
		TargetRolesAccountResult.Decode(&targetRolesAccount)

		if userRolesAccount.Name == "Admin" {
			if targetRolesAccount.Name != "Facility" {
				c.JSON(http.StatusBadRequest, bson.M{
					"Status":  http.StatusBadRequest,
					"Message": "Anda tidak bisa mendaftarkan akun selain akun dengan jenis akun Fasilitas",
				})
				return
			}
		}
		if userRolesAccount.Name == "Facility" {
			if !(targetRolesAccount.Name == "Mother" || targetRolesAccount.Name == "Child") {
				c.JSON(http.StatusBadRequest, bson.M{
					"Status":  http.StatusBadRequest,
					"Message": "Anda tidak bisa mendaftarkan akun selain akun dengan jenis akun Ibu atau Anak",
				})
				return
			}
		}

		bytes, errors := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if errors != nil {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Password tidak valid"}})
		}

		rolesCount, rolesGetError = rolesCollection.CountDocuments(ctx, bson.M{"_id": user.RolesId})

		if rolesGetError != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Internal server error"}})
			return
		}

		if rolesCount == 0 {
			c.JSON(http.StatusBadRequest, views.MasterResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Roles ID not configured"}})
			return
		}

		newUser := models.User{
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Identifier: user.Identifier,
			Address:    user.Address,
			RolesId:    user.RolesId,
			ParentId:   accountUserId,
			Password:   string(bytes),
			UpdatedAt:  time.Now().UTC(),
			CreatedAt:  time.Now().UTC(),
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
				"id":         result.InsertedID,
				"firstname":  newUser.FirstName,
				"lastname":   newUser.LastName,
				"roles":      user.RolesId,
				"identifier": newUser.Identifier,
			},
		})
	}
}
