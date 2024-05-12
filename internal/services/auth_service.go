package services

import (
	"net/http"
	"social-network-api/internal/dtos"
	"social-network-api/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(request dtos.RegisterRequest) *utils.ApiResponse[dtos.LoginResponse] {

	return utils.NewApiResponse[dtos.LoginResponse](http.StatusOK, "hola", nil)
}
