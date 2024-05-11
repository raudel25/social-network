package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct{}

func (uc *UserController) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "¡Hola Estoy haciendo inicio de sesión.",
	})
}

func (uc *UserController) Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario registrado exitosamente.",
	})
}

func NewUserController() *UserController {
	return &UserController{}
}
