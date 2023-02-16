package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	mongoUri := ""
	if gin.Mode() == "debug" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		mongoUri = os.Getenv("MONGOURI")
	} else {
		mongoUri = os.Getenv("MONGOURI")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("NoStunting").Collection(collectionName)
	return collection
}
