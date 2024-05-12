package controllers

import (
	"social-network-api/internal/dtos"
	"social-network-api/internal/services"
	"social-network-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	authService *services.AuthService
}

func (s *UserController) Login(c *gin.Context) {
	var request dtos.LoginRequest

	checkRequest := utils.CheckRequest(c, &request)
	if !checkRequest.Ok() {
		checkRequest.Response(c)
		return
	}

	s.authService.Login(request).Response(c)

}

func (s *UserController) Register(c *gin.Context) {
	var request dtos.RegisterRequest

	checkRequest := utils.CheckRequest(c, &request)
	if !checkRequest.Ok() {
		checkRequest.Response(c)
		return
	}

	s.authService.Register(request).Response(c)
}

func NewUserController(service *services.AuthService) *UserController {
	return &UserController{authService: service}
}
