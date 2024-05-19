package services

import (
	"social-network-api/internal/models"
	"social-network-api/internal/pkg"

	"gorm.io/gorm"
)

type PhotoService struct {
	db *gorm.DB
}

func NewPhotoService(db *gorm.DB) *PhotoService {
	return &PhotoService{db: db}
}

func (s *PhotoService) GetPhoto(id uint) *pkg.ApiResponse[models.Photo] {
	var photo models.Photo
	if s.db.First(&photo, id).Error != nil {
		return pkg.NewNotFound[models.Photo]("Photo not found")
	}

	return pkg.NewOk(&photo)
}

func (s *PhotoService) UploadPhoto(filename string) *pkg.ApiResponse[uint] {
	photo := models.Photo{Filename: filename}
	s.db.Create(&photo)

	return pkg.NewOk(&photo.ID)
}
