package controllers

import (
	"net/http"
	"social-network-api/internal/core"
	"social-network-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckRequest[T any](c *gin.Context, request *T) *core.ApiResponse[T] {
	if err := c.ShouldBindJSON(&request); err != nil {
		return core.NewBadRequest[T](err.Error())
	}

	return core.NewOk[T]()
}

func CheckAuthorized(c *gin.Context, jwt_service *services.JWTService) *core.ApiResponse[core.JWTDto] {
	token, err := jwt_service.CheckJWT(c.GetHeader("Authorization"))

	if err != nil {
		return core.NewApiResponse[core.JWTDto](http.StatusUnauthorized, "Invalid token", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return core.NewApiResponse[core.JWTDto](http.StatusUnauthorized, "Invalid token", nil)
	}

	return core.NewApiResponse[core.JWTDto](http.StatusOK, "", &core.JWTDto{Id: uint(claims["id"].(float64)), Username: claims["username"].(string)})
}
