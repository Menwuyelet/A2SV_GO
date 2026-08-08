package main

import (
	"task_manager/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routers.StartRoute(router)

	router.Run()
}
