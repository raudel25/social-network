package utils

import (
	"github.com/gin-gonic/gin"
)

func CheckRequest[T any](c *gin.Context, request *T) *ApiResponse[T] {
	if err := c.ShouldBindJSON(&request); err != nil {
		return NewBadRequest[T](err.Error())
	}

	return NewOk[T]()
}
