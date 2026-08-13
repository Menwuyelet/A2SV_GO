package controllers

import (
	"errors"
	"net/http"
	"task_manager/Domain/dto"
	"task_manager/Usecases"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	service Usecases.TaskUsecase
}

func NewTaskController(service Usecases.TaskUsecase) *TaskController {
	return &TaskController{service: service}
}

func (tc *TaskController) AddTask(ctx *gin.Context) {
	var newTask dto.TaskRequest

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	owner, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated user."})
		return
	}

	task, err := tc.service.AddTask(ctx.Request.Context(), newTask, owner.(string))

	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task."})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, task)
}

func (tc *TaskController) ListAllTasks(ctx *gin.Context) {
	owner, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated user."})
		return
	}

	tasks, err := tc.service.ListTasks(ctx.Request.Context(), owner.(string))
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (tc *TaskController) RetrieveTask(ctx *gin.Context) {
	id := ctx.Param("id")
	ownerID, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized user."})
		return
	}

	task, err := tc.service.GetTask(ctx.Request.Context(), id, ownerID.(string))
	if err != nil {
		if errors.Is(err, Usecases.ErrUnauthorized) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized user."})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func (tc *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	owner, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated user."})
		return
	}

	var updatedTask dto.TaskUpdateRequest

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request."})
		return
	}

	task, err := tc.service.UpdateTask(ctx.Request.Context(), id, owner.(string), updatedTask)

	if err != nil {
		if errors.Is(err, Usecases.ErrUnauthorized) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user."})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	owner, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated user."})
		return
	}

	err := tc.service.DeleteTask(ctx, id, owner.(string))
	if err != nil {
		if errors.Is(err, Usecases.ErrUnauthorized) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user."})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
