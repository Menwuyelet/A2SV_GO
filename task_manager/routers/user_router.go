package routers

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, uc *controllers.UserController) {

	// user
	rg.POST("/user/register", uc.Register)
	rg.POST("/user/login", uc.Login)

}
