package controllers

import (
	"net/http"
	"task_manager/models"
	"task_manager/repository"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	repo repository.TaskRepository
}

func NewTaskController(repo repository.TaskRepository) *TaskController {
	return &TaskController{repo: repo}
}

func (tc *TaskController) AddTask(ctx *gin.Context) {
	var newTask models.Task

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := tc.repo.AddTask(ctx.Request.Context(), newTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task."})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, task)
}

func (tc *TaskController) ListAllTasks(ctx *gin.Context) {
	tasks, err := tc.repo.ListTasks(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (tc *TaskController) RetrieveTask(ctx *gin.Context) {
	id := ctx.Param("id")

	task, err := tc.repo.GetTask(ctx.Request.Context(), id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func (tc *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request."})
		return
	}

	task, err := tc.repo.UpdateTask(ctx.Request.Context(), id, updatedTask)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := tc.repo.DeleteTask(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Task deleted successfully."})
}
