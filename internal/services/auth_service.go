package services

import (
	"net/http"
	"regexp"
	"social-network-api/internal/dtos"
	"social-network-api/internal/models"
	"social-network-api/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // El segundo parámetro es el costo, ajustable según necesidades
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) Register(request dtos.RegisterRequest) *utils.ApiResponse[dtos.LoginResponse] {

	if s.db.Where("username =?", request.Username).First(&models.User{}).Error == nil || s.db.Where("email =?", request.Email).First(&models.User{}).Error == nil {
		return utils.NewApiResponse[dtos.LoginResponse](http.StatusConflict, "User already exists", nil)
	}

	if len(request.Password) < 8 {
		return utils.NewApiResponse[dtos.LoginResponse](http.StatusBadRequest, "Password must be at least 8 characters", nil)
	}

	if !isEmailValid(request.Email) {
		return utils.NewApiResponse[dtos.LoginResponse](http.StatusBadRequest, "Invalid email", nil)
	}

	hashedPassword, _ := HashPassword(request.Password)
	s.db.Create(&models.User{Name: request.Name, Username: request.Username, Email: request.Email, Password: hashedPassword})

	return utils.NewApiResponse(http.StatusOK, "Ok", &dtos.LoginResponse{Name: request.Name, Username: request.Username, Email: request.Email, Token: ""})
}

func (s *AuthService) Login(request dtos.LoginRequest) *utils.ApiResponse[dtos.LoginResponse] {

	var user models.User
	if s.db.Where("username = ? OR email = ?", request.Username, request.Username).First(&user).Error != nil {
		return utils.NewApiResponse[dtos.LoginResponse](http.StatusNotFound, "Incorrect user or password", nil)
	}

	if !CheckPasswordHash(request.Password, user.Password) {
		return utils.NewApiResponse[dtos.LoginResponse](http.StatusUnauthorized, "Incorrect user or password", nil)
	}

	return utils.NewApiResponse(http.StatusOK, "Ok", &dtos.LoginResponse{Name: user.Name, Username: user.Username, Email: user.Email, Token: ""})
}
