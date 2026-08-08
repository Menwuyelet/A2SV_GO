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

	status := data.AddTask(newTask)

	if status != "Success" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request."})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newTask)
}

func ListAllTasks(ctx *gin.Context) {
	tasks := data.ListTasks()

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func RetrieveTask(ctx *gin.Context) {
	id := ctx.Param("id")

	task, status := data.GetTask(id)

	if status != nil {
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

	task, status := data.UpdateTask(id, updatedTask)

	if status != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	status := data.DeleteTask(id)

	if status != "Success" {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Tasks deleted."})
}
