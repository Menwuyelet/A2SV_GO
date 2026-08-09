package controllers

import (
	"net/http"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func AddTask(ctx *gin.Context) {
	var newTask models.Task

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := data.AddTask(newTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task."})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, task)
}

func ListAllTasks(ctx *gin.Context) {
	tasks, err := data.ListTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func RetrieveTask(ctx *gin.Context) {
	id := ctx.Param("id")

	task, err := data.GetTask(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request."})
		return
	}

	task, err := data.UpdateTask(id, updatedTask)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := data.DeleteTask(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Task deleted successfully."})
}
