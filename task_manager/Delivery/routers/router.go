package routers

import (
	"task_manager/Delivery/controllers"
	"task_manager/Infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Task *controllers.TaskController
	User *controllers.UserController
}

func StartRoute(c *Controllers) *gin.Engine {
	router := gin.Default()

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	RegisterTaskRoutes(protected, c.Task)

	public := router.Group("/api")
	RegisterUserRoutes(public, c.User)
	return router
}
