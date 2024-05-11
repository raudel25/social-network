package config

import (
	"social-network-api/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	authRoutes(r)

	return r
}

func authRoutes(r *gin.Engine) {
	controller := controllers.NewUserController()
	auth := r.Group("/auth")

	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Login)
}
