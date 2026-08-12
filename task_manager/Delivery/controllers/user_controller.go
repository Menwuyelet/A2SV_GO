package controllers

import (
	"net/http"
	"task_manager/Domain/dto"
	"task_manager/Infrastructure/utils"
	"task_manager/Usecases"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *Usecases.UserService
}

func NewUserController(service *Usecases.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) Register(ctx *gin.Context) {
	var newUser dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		utils.HandleValidationError(ctx, err)
		return
	}

	user, err := uc.service.Register(ctx.Request.Context(), newUser)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (uc *UserController) Login(ctx *gin.Context) {
	var newLogin dto.LoginRequest

	if err := ctx.ShouldBindJSON(&newLogin); err != nil {
		utils.HandleValidationError(ctx, err)
		return
	}

	token, err := uc.service.Login(ctx.Request.Context(), newLogin)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
