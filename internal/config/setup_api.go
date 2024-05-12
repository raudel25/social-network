package config

import (
	"social-network-api/internal/controllers"
	"social-network-api/internal/services"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupApi(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	authService := services.NewAuthService(db)

	authRoutes(r, authService)

	return r
}

func authRoutes(r *gin.Engine, service *services.AuthService) {
	controller := controllers.NewUserController(service)
	auth := r.Group("/auth")

	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Register)
}
