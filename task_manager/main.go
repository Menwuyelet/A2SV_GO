package main

import (
	"context"
	"log"
	"os"

	"task_manager/Delivery/controllers"
	"task_manager/Delivery/routers"
	"task_manager/Usecases"
	"task_manager/repository"

	"github.com/joho/godotenv"
)

func setupTaskController() *controllers.TaskController {
	repo := repository.NewMongoTaskRepository(repository.GetTaskCollection())
	service := Usecases.NewTaskService(repo)
	return controllers.NewTaskController(service)
}

func setupUserController() *controllers.UserController {
	repo := repository.NewMongoUserRepository(repository.GetUserCollection())
	service := Usecases.NewUserService(repo)
	return controllers.NewUserController(service)
}

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

	// taskRepo := repository.NewMongoTaskRepository(repository.GetTaskCollection())

	// taskService := Usecases.NewTaskService(taskRepo)

	// taskController := controllers.NewTaskController(taskService)

	// userRepo := repository.NewMongoUserRepository(repository.GetUserCollection())

	// userService := Usecases.NewUserService(userRepo)

	// userController := controllers.NewUserController(userService)

	allControllers := &routers.Controllers{
		Task: setupTaskController(),
		User: setupUserController(),
	}

	// handles routers

	router := routers.StartRoute(allControllers)

	router.Run()
}
