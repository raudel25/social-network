package controllers

import (
	"net/http"
	"social-network-api/internal/models"
	"social-network-api/internal/pkg"
	"social-network-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckRequest[T any](c *gin.Context, request *T) *pkg.ApiResponse[T] {
	if err := c.ShouldBindJSON(&request); err != nil {
		return pkg.NewBadRequest[T](err.Error())
	}

	return pkg.NewOk[T](nil)
}

func CheckAuthorized(c *gin.Context, jwt_service *services.JWTService) *pkg.ApiResponse[models.JWTDto] {
	token, err := jwt_service.CheckJWT(c.GetHeader("Authorization"))

	if err != nil {
		return pkg.NewApiResponse[models.JWTDto](http.StatusUnauthorized, "Invalid token", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return pkg.NewApiResponse[models.JWTDto](http.StatusUnauthorized, "Invalid token", nil)
	}

	return pkg.NewApiResponse(http.StatusOK, "", &models.JWTDto{ID: uint(claims["id"].(float64)), Username: claims["username"].(string)})
}

func CheckId(c *gin.Context) *pkg.ApiResponse[uint] {
	id := c.Param("id")
	num, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return pkg.NewBadRequest[uint]("Incorrect url")
	}

	finalUint := uint(num)
	return pkg.NewOk(&finalUint)
}

func CheckPagination[T any](c *gin.Context) *pkg.ApiResponse[pkg.Pagination[T]] {
	page, err := strconv.ParseUint(c.Query("page"), 10, 64)
	if err != nil {
		return pkg.NewBadRequest[pkg.Pagination[T]]("Incorrect page")
	}

	limit, err := strconv.ParseUint(c.Query("limit"), 10, 64)
	if err != nil {
		return pkg.NewBadRequest[pkg.Pagination[T]]("Incorrect limit")
	}

	return pkg.NewOk(&pkg.Pagination[T]{Page: int(page), Limit: int(limit)})
}
