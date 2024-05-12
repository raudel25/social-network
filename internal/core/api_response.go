package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse[T any] struct {
	Status  int
	Message string
	Data    *T
}

func (uc *ApiResponse[T]) Ok() bool {
	return uc.Status == http.StatusOK
}

func (uc *ApiResponse[T]) Response(c *gin.Context) {
	if uc.Data == nil {
		c.JSON(uc.Status, gin.H{
			"message": uc.Message,
		})

		return
	}

	c.JSON(uc.Status, gin.H{
		"message": uc.Message,
		"data":    uc.Data})

}

func NewApiResponse[T any](status int, message string, data *T) *ApiResponse[T] {
	return &ApiResponse[T]{Status: status, Message: message, Data: data}
}

func NewBadRequest[T any](message string) *ApiResponse[T] {
	return &ApiResponse[T]{Status: http.StatusBadRequest, Message: message, Data: nil}
}

func NewOk[T any]() *ApiResponse[T] {
	return &ApiResponse[T]{Status: http.StatusOK, Message: "", Data: nil}
}
