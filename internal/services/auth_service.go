package services

import (
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"social-network-api/internal/core"
	"social-network-api/internal/models"
)

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // El segundo parámetro es el costo, ajustable según necesidades
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type AuthService struct {
	db         *gorm.DB
	jwtService *JWTService
}

func NewAuthService(db *gorm.DB, jwtService *JWTService) *AuthService {
	return &AuthService{db: db, jwtService: jwtService}
}

func (s *AuthService) Register(request core.RegisterRequest) *core.ApiResponse[core.LoginResponse] {

	if s.db.Where("username =?", request.Username).First(&models.User{}).Error == nil || s.db.Where("email =?", request.Email).First(&models.User{}).Error == nil {
		return core.NewApiResponse[core.LoginResponse](http.StatusConflict, "User already exists", nil)
	}

	if len(request.Password) < 8 {
		return core.NewApiResponse[core.LoginResponse](http.StatusBadRequest, "Password must be at least 8 characters", nil)
	}

	if !isEmailValid(request.Email) {
		return core.NewApiResponse[core.LoginResponse](http.StatusBadRequest, "Invalid email", nil)
	}

	hashedPassword, _ := HashPassword(request.Password)
	user := models.User{Name: request.Name, Username: request.Username, Email: request.Email, Password: hashedPassword}
	s.db.Create(&user)

	token, _ := s.jwtService.GenerateJWT(user.ID, user.Username)
	return core.NewApiResponse(http.StatusOK, "Ok", &core.LoginResponse{Name: request.Name, Username: request.Username, Email: request.Email, Token: token})
}

func (s *AuthService) Login(request core.LoginRequest) *core.ApiResponse[core.LoginResponse] {

	var user models.User
	if s.db.Where("username = ? OR email = ?", request.Username, request.Username).First(&user).Error != nil {
		return core.NewApiResponse[core.LoginResponse](http.StatusNotFound, "Incorrect user or password", nil)
	}

	if !CheckPasswordHash(request.Password, user.Password) {
		return core.NewApiResponse[core.LoginResponse](http.StatusUnauthorized, "Incorrect user or password", nil)
	}

	token, _ := s.jwtService.GenerateJWT(user.ID, user.Username)
	return core.NewApiResponse(http.StatusOK, "Ok", &core.LoginResponse{Name: user.Name, Username: user.Username, Email: user.Email, Token: token})
}

func (s *AuthService) Renew(request *core.JWTDto) *core.ApiResponse[core.LoginResponse] {
	user := models.User{}
	s.db.First(&user, request.Id)

	token, _ := s.jwtService.GenerateJWT(user.ID, user.Username)
	return core.NewApiResponse(http.StatusOK, "Ok", &core.LoginResponse{Name: user.Name, Username: user.Username, Email: user.Email, Token: token})
}
