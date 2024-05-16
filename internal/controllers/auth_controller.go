package controllers

import (
	"social-network-api/internal/models"
	"social-network-api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	authService *services.AuthService
	jwtService  *services.JWTService
}

func (s *UserController) Login(c *gin.Context) {
	var request models.LoginRequest

	checkRequest := CheckRequest(c, &request)
	if !checkRequest.Ok() {
		checkRequest.Response(c)
		return
	}

	s.authService.Login(&request).Response(c)

}

func (s *UserController) Register(c *gin.Context) {
	var request models.RegisterRequest

	checkRequest := CheckRequest(c, &request)
	if !checkRequest.Ok() {
		checkRequest.Response(c)
		return
	}

	s.authService.Register(&request).Response(c)
}

func (s *UserController) Renew(c *gin.Context) {
	checkAuthorized := CheckAuthorized(c, s.jwtService)

	if !checkAuthorized.Ok() {
		checkAuthorized.Response(c)
		return
	}

	s.authService.Renew(checkAuthorized.Data).Response(c)

}

func NewUserController(authService *services.AuthService, jwtService *services.JWTService) *UserController {
	return &UserController{authService: authService, jwtService: jwtService}
}
