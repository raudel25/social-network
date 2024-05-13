package controllers

import (
	"net/http"
	"social-network-api/internal/models"
	"social-network-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckRequest[T any](c *gin.Context, request *T) *models.ApiResponse[T] {
	if err := c.ShouldBindJSON(&request); err != nil {
		return models.NewBadRequest[T](err.Error())
	}

	return models.NewOk[T](nil)
}

func CheckAuthorized(c *gin.Context, jwt_service *services.JWTService) *models.ApiResponse[models.JWTDto] {
	token, err := jwt_service.CheckJWT(c.GetHeader("Authorization"))

	if err != nil {
		return models.NewApiResponse[models.JWTDto](http.StatusUnauthorized, "Invalid token", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.NewApiResponse[models.JWTDto](http.StatusUnauthorized, "Invalid token", nil)
	}

	return models.NewApiResponse[models.JWTDto](http.StatusOK, "", &models.JWTDto{ID: uint(claims["id"].(float64)), Username: claims["username"].(string)})
}

func CheckId(c *gin.Context) *models.ApiResponse[uint] {
	id := c.Param("id")
	num, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return models.NewBadRequest[uint]("Incorrect url")
	}

	finalUint := uint(num)
	return models.NewOk[uint](&finalUint)
}
