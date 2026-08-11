package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// handles db connection
var DB *mongo.Database

func DBconnect(mongoURI string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")

	DB = client.Database("task_management_api")

	return client, nil
}

// handles collection accessing
func GetTaskCollection() *mongo.Collection {
	return DB.Collection("task")
}

func GetUserCollection() *mongo.Collection {
	return DB.Collection("user")
}
