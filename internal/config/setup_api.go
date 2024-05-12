package config

import (
	"social-network-api/internal/controllers"
	"social-network-api/internal/db"
	"social-network-api/internal/models"
	"social-network-api/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupApi(config models.Config) *gin.Engine {
	r := gin.Default()

	db := db.ConnectDatabase(config)

	jwtService := services.NewJwtService(config.SecretKey)
	authService := services.NewAuthService(db, jwtService)

	authRoutes(r, authService)

	return r
}

func authRoutes(r *gin.Engine, service *services.AuthService) {
	controller := controllers.NewUserController(service)
	auth := r.Group("/auth")

	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Register)
}
