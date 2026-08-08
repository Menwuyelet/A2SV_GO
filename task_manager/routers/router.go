package routers

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func StartRoute(router *gin.Engine) {

	router.POST("/tasks", controllers.AddTask)
	router.GET("/tasks", controllers.ListAllTasks)
	router.GET("/tasks/:id", controllers.RetrieveTask)
	router.PUT("/tasks/:id", controllers.UpdateTask)
	router.DELETE("/tasks/:id", controllers.DeleteTask)

}
