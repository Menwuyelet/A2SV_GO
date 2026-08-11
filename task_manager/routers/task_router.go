package routers

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(rg *gin.RouterGroup, tc *controllers.TaskController) {

	rg.POST("/tasks", tc.AddTask)
	rg.GET("/tasks", tc.ListAllTasks)
	rg.GET("/tasks/:id", tc.RetrieveTask)
	rg.PUT("/tasks/:id", tc.UpdateTask)
	rg.DELETE("/tasks/:id", tc.DeleteTask)

}
