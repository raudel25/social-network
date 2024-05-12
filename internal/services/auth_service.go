package services

import (
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

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

func (s *AuthService) Register(request models.RegisterRequest) *models.ApiResponse[models.LoginResponse] {

	if s.db.Where("username =?", request.Username).First(&models.User{}).Error == nil || s.db.Where("email =?", request.Email).First(&models.User{}).Error == nil {
		return models.NewApiResponse[models.LoginResponse](http.StatusConflict, "User already exists", nil)
	}

	if len(request.Password) < 8 {
		return models.NewApiResponse[models.LoginResponse](http.StatusBadRequest, "Password must be at least 8 characters", nil)
	}

	if !isEmailValid(request.Email) {
		return models.NewApiResponse[models.LoginResponse](http.StatusBadRequest, "Invalid email", nil)
	}

	hashedPassword, _ := HashPassword(request.Password)
	user := models.User{Username: request.Username, Email: request.Email, Password: hashedPassword}
	s.db.Create(&user)

	profile := models.Profile{UserID: user.ID, Name: request.Name, ProfilePhotoID: request.ProfilePhotoID, BannerPhotoID: request.BannerPhotoID, RichText: request.RichText}
	s.db.Create(&profile)

	token, _ := s.jwtService.GenerateJWT(profile.ID, user.Username)
	return models.NewApiResponse(http.StatusOK, "Ok", &models.LoginResponse{Username: request.Username, Token: token, Profile: profile})
}

func (s *AuthService) Login(request models.LoginRequest) *models.ApiResponse[models.LoginResponse] {

	var user models.User
	if s.db.Where("username = ? OR email = ?", request.Username, request.Username).First(&user).Error != nil {
		return models.NewApiResponse[models.LoginResponse](http.StatusNotFound, "Incorrect user or password", nil)
	}

	if !CheckPasswordHash(request.Password, user.Password) {
		return models.NewApiResponse[models.LoginResponse](http.StatusUnauthorized, "Incorrect user or password", nil)
	}

	var profile models.Profile
	s.db.Preload("Follows").Preload("FollowedBy").Where("user_id = ?", user.ID).First(&profile)

	token, _ := s.jwtService.GenerateJWT(profile.ID, user.Username)
	return models.NewApiResponse(http.StatusOK, "Ok", &models.LoginResponse{Username: user.Username, Token: token, Profile: profile})
}

func (s *AuthService) Renew(request *models.JWTDto) *models.ApiResponse[models.LoginResponse] {
	profile := models.Profile{}
	s.db.Preload("Follows").Preload("FollowedBy").Preload("User").First(&profile, request.Id)

	token, _ := s.jwtService.GenerateJWT(profile.ID, profile.User.Username)
	return models.NewApiResponse(http.StatusOK, "Ok", &models.LoginResponse{Username: profile.User.Username, Token: token, Profile: profile})
}
