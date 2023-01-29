package controllers

import (
	"be-no-stunting-v2/configs"
	"be-no-stunting-v2/views"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var rolesCollectionOnMaster *mongo.Collection = configs.GetCollection(configs.DB, "roles")

func GetMasterData() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var roles []views.RolesWithId
		defer cancel()

		rolesResult, err := rolesCollectionOnMaster.Find(ctx, bson.M{})
		rolesResult.Decode(&roles)

		if err != nil {
			c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer rolesResult.Close(ctx)

		for rolesResult.Next(ctx) {
			var singleRoles views.RolesWithId
			if err = rolesResult.Decode(&singleRoles); err != nil {
				c.JSON(http.StatusInternalServerError, views.MasterResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			roles = append(roles, singleRoles)
		}

		// values := bson.M{
		// 	"data1": 88,
		// 	"data2": 77,
		// }

		// json_data, err__ := json.Marshal(values)
		// if err__ != nil {
		// 	log.Fatal("FUCJK")
		// }

		// testRequest, _ := http.Post("https://gdsc-pens-iot-listener-lxz6xwlfka-et.a.run.app/no-stunting/8890", "application/json", bytes.NewBuffer(json_data))
		// fmt.Println(testRequest)

		c.JSON(http.StatusOK,
			bson.M{
				"Status":  http.StatusOK,
				"Message": "Success",
				"Data": bson.M{
					"Roles": roles,
				},
			},
		)
	}
}
