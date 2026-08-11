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

	taskRepo := repository.NewMongoTaskRepository(repository.GetTaskCollection())

	taskService := data.NewTaskService(taskRepo)

	taskController := controllers.NewTaskController(*taskService)

	userRepo := repository.NewMongoUserRepository(repository.GetUserCollection())

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
