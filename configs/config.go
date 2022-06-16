package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection

func ConnectDB() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	//client options
	clientOptions := options.Client().ApplyURI(os.Getenv("connectionstring"))
	//connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDb Connection success")
	Collection = client.Database("Notes_Data").Collection("notes")
	//collection instance
	fmt.Println("Collection instance is ready")
	data, err := Collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data of collection", data)
}
