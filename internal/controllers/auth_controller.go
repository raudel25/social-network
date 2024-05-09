package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
        "message": "¡Hola! Estoy haciendo inicio de sesión.",
    })
}

func Register(c *gin.Context) {
}
