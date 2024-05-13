package services

import (
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"social-network-api/internal/models"
	"social-network-api/internal/pkg"
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

func profileToResponseAuth(profile *models.Profile, username string) *models.ProfileResponse {
	return &models.ProfileResponse{
		ID:           profile.ID,
		Name:         profile.Name,
		ProfilePhoto: profile.ProfilePhoto,
		BannerPhoto:  profile.BannerPhoto,
		RichText:     profile.RichText,
		Follow:       false,
		Username:     username,
	}
}

type AuthService struct {
	db         *gorm.DB
	jwtService *JWTService
}

func NewAuthService(db *gorm.DB, jwtService *JWTService) *AuthService {
	return &AuthService{db: db, jwtService: jwtService}
}

func (s *AuthService) Register(request models.RegisterRequest) *pkg.ApiResponse[models.LoginResponse] {
	if s.db.Where("username =?", request.Username).First(&models.User{}).Error == nil || s.db.Where("email =?", request.Email).First(&models.User{}).Error == nil {
		return pkg.NewApiResponse[models.LoginResponse](http.StatusConflict, "User already exists", nil)
	}

	if len(request.Password) < 8 {
		return pkg.NewApiResponse[models.LoginResponse](http.StatusBadRequest, "Password must be at least 8 characters", nil)
	}

	if !isEmailValid(request.Email) {
		return pkg.NewApiResponse[models.LoginResponse](http.StatusBadRequest, "Invalid email", nil)
	}

	if request.Profile.ProfilePhotoID != nil && s.db.First(&models.Photo{}, request.Profile.ProfilePhotoID).Error != nil {
		return pkg.NewNotFound[models.LoginResponse]("Profile photo not found")
	}

	if request.Profile.BannerPhotoID != nil && s.db.First(&models.Photo{}, request.Profile.BannerPhotoID).Error != nil {
		return pkg.NewNotFound[models.LoginResponse]("Banner photo not found")
	}

	hashedPassword, _ := HashPassword(request.Password)
	user := models.User{Username: request.Username, Email: request.Email, Password: hashedPassword}
	s.db.Create(&user)

	profile := models.Profile{UserID: user.ID, Name: request.Name, ProfilePhotoID: request.Profile.ProfilePhotoID, BannerPhotoID: request.Profile.BannerPhotoID, RichText: request.Profile.RichText}
	s.db.Create(&profile)

	token, _ := s.jwtService.GenerateJWT(profile.ID, user.Username)
	return pkg.NewOk(&models.LoginResponse{Username: request.Username, Token: token, Profile: *profileToResponseAuth(&profile, user.Username)})
}

func (s *AuthService) Login(request models.LoginRequest) *pkg.ApiResponse[models.LoginResponse] {

	var user models.User
	if s.db.Where("username = ? OR email = ?", request.Username, request.Username).First(&user).Error != nil {
		return pkg.NewApiResponse[models.LoginResponse](http.StatusNotFound, "Incorrect user or password", nil)
	}

	if !CheckPasswordHash(request.Password, user.Password) {
		return pkg.NewApiResponse[models.LoginResponse](http.StatusUnauthorized, "Incorrect user or password", nil)
	}

	var profile models.Profile
	s.db.Where("user_id = ?", user.ID).Preload("ProfilePhoto").Preload("BannerPhoto").First(&profile)

	token, _ := s.jwtService.GenerateJWT(profile.ID, user.Username)

	return pkg.NewOk(&models.LoginResponse{Username: user.Username, Token: token, Profile: *profileToResponseAuth(&profile, user.Username)})
}

func (s *AuthService) Renew(request *models.JWTDto) *pkg.ApiResponse[models.LoginResponse] {
	profile := models.Profile{}
	s.db.Preload("User").Preload("ProfilePhoto").Preload("BannerPhoto").First(&profile, request.ID)

	token, _ := s.jwtService.GenerateJWT(profile.ID, profile.User.Username)
	return pkg.NewOk(&models.LoginResponse{Username: profile.User.Username, Token: token, Profile: *profileToResponseAuth(&profile, profile.User.Username)})
}
