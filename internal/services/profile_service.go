package services

import (
	"social-network-api/internal/models"

	"gorm.io/gorm"
)

type ProfileService struct {
	db *gorm.DB
}

func (s *ProfileService) EditProfile(request *models.ProfileDto, jwt *models.JWTDto) *models.SingleApiResponse {
	if request.ProfilePhotoID != nil && s.db.First(&models.Photo{}, request.ProfilePhotoID).Error != nil {
		return models.NewSingleNotFound("Profile photo not found")
	}

	if request.BannerPhotoID != nil && s.db.First(&models.Photo{}, request.BannerPhotoID).Error != nil {
		return models.NewSingleNotFound("Banner photo not found")
	}

	var profile models.Profile
	if s.db.Find(&profile, jwt.ID).Error != nil {
		return models.NewSingleNotFound("Profile not found")
	}

	profile.Name = request.Name
	profile.ProfilePhotoID = request.ProfilePhotoID
	profile.BannerPhotoID = request.BannerPhotoID
	profile.RichText = request.RichText

	s.db.Where("id =?", jwt.ID).Updates(&profile)

	return models.NewSingleOkSingle()
}

func (s *ProfileService) FollowUnFollow(id uint, jwt *models.JWTDto) *models.SingleApiResponse {
	if s.db.First(&models.Profile{}, id).Error != nil {
		return models.NewSingleNotFound("Profile not found")
	}

	if s.db.Where("follower_id =? AND followed_id =?", jwt.ID, id).First(&models.Follow{}).Error != nil {
		s.db.Create(&models.Follow{FollowerProfileID: jwt.ID, FollowedProfileID: id})
	} else {
		s.db.Where("follower_id =? AND followed_id =?", jwt.ID, id).Delete(&models.Follow{})
	}

	return models.NewSingleOkSingle()
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{db: db}
}
