package routes

import (
	"social-network-api/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/login", controllers.Login)

	return r
}
