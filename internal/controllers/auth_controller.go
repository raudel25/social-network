package controllers

import (
	"net/http"

	"social-network-api/internal/dtos"
	"social-network-api/internal/services"
	"social-network-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	authService *services.AuthService
}

func (uc *UserController) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "¡Hola Estoy haciendo inicio de sesión.",
	})
}

func (uc *UserController) Register(c *gin.Context) {
	var request dtos.RegisterRequest

	checkRequest := utils.CheckRequest(c, &request)
	if checkRequest.Ok() {
		checkRequest.Response(c)
		return
	}

	uc.authService.Register(request).Response(c)
}

func NewUserController(service *services.AuthService) *UserController {
	return &UserController{authService: service}
}
