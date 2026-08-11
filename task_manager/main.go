package main

import (
	"context"
	"log"
	"os"

	"task_manager/controllers"
	"task_manager/data"
	"task_manager/repository"
	"task_manager/routers"

	"github.com/joho/godotenv"
)

func main() {
	// loads .env
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	mongoURI := os.Getenv("MONGO_URI")

	// creates db connection by passing the URI from .env
	client, err := repository.DBconnect(mongoURI)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	taskRepo := data.NewMongoTaskRepository(repository.GetTaskCollection())

	taskController := controllers.NewTaskController(taskRepo)

	userRepo := repository.NewMongoUserRepository(repository.GetTaskCollection())

	userService := data.NewUserService(userRepo)

	userController := controllers.NewUserController(userService)

	allControllers := &routers.Controllers{
		Task: taskController,
		User: userController,
	}

	// handles routers

	router := routers.StartRoute(allControllers)

	router.Run()
}
