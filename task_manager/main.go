package main

import (
	"context"
	"log"
	"os"
	"task_manager/data"
	"task_manager/models"
	"task_manager/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	// loads .env
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	mongoURI := os.Getenv("MONGO_URI")

	// creates db connection by passing the URI from .env
	client, err := models.DBconnect(mongoURI)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// retrieves the collection we are accessing
	collection := models.GetCollection("tasks")

	// initiates the collection for crud operations in data module
	data.InitCollection(collection)

	seedTasks(collection)

	// handles routers
	router := gin.Default()

	routers.StartRoute(router)

	router.Run()
}

func seedTasks(collection *mongo.Collection) {
	// checks if the system have any existing data and if so stop the function and return
	count, err := collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Failed to count tasks: %v", err)
		return
	}

	if count > 0 {
		return
	}

	// else it will seed 3 new tasks for testing
	tasks := []models.Task{
		{Title: "Task 1", Description: "First task", DueDate: "2026-10-11", Status: "Pending"},
		{Title: "Task 2", Description: "Second task", DueDate: "2026-10-11", Status: "In Progress"},
		{Title: "Task 3", Description: "Third task", DueDate: "2026-10-11", Status: "Completed"},
	}

	for _, task := range tasks {
		task.ID = bson.NewObjectID()
		if _, err := collection.InsertOne(context.TODO(), task); err != nil {
			log.Printf("Failed to seed task: %v", err)
		}
	}
}
