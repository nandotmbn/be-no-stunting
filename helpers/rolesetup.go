package helpers

import (
	"be-no-stunting-v2/configs"
	"be-no-stunting-v2/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var rolesCollection *mongo.Collection = configs.GetCollection(configs.DB, "roles")

func roleSetup(rolesName string, displayName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"$or": []bson.M{ // you can try this in []interface
		{"name": rolesName},
		{"displayname": displayName},
	}}

	count, err := rolesCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("Terdapat kesalahan perhitungan dokumen setup role")
		return
	}

	if count >= 1 {
		fmt.Println("Role dengan nama ini sudah pernah didaftarkan")
		return
	}

	newRoles := models.Roles{
		Name:        rolesName,
		DisplayName: displayName,
	}

	result, err := rolesCollection.InsertOne(ctx, newRoles)
	if err != nil {
		fmt.Println("Terdapat kesalahan perhitungan dokumen setup role")
		return
	}

	fmt.Println(result.InsertedID)
}

func RolesSetup() {
	roleSetup("Child", "Anak")
	roleSetup("Facility", "Fasilitas")
	roleSetup("Mother", "Ibu")
	roleSetup("Admin", "Administrator")
}
