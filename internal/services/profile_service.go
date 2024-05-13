package services

import (
	"net/http"
	"social-network-api/internal/models"

	"gorm.io/gorm"
)

type ProfileService struct {
	db *gorm.DB
}

func (s *ProfileService) EditProfile(request *models.ProfileDto, id uint, jwt *models.JWTDto) *models.SingleApiResponse {
	if request.ProfilePhotoID != nil && s.db.First(&models.Photo{}, request.ProfilePhotoID).Error != nil {
		return models.NewSingleNotFound("Profile photo not found")
	}

	if request.BannerPhotoID != nil && s.db.First(&models.Photo{}, request.BannerPhotoID).Error != nil {
		return models.NewSingleNotFound("Banner photo not found")
	}

	if id != jwt.ID {
		return models.NewSingleApiResponse(http.StatusUnauthorized, "Unauthorized")
	}

	var profile models.Profile
	if s.db.Find(&profile, id).Error != nil {
		return models.NewSingleNotFound("Profile not found")
	}

	profile.Name = request.Name
	profile.ProfilePhotoID = request.ProfilePhotoID
	profile.BannerPhotoID = request.BannerPhotoID
	profile.RichText = request.RichText

	s.db.Where("id =?", id).Updates(&profile)

	return models.NewSingleOkSingle()
}

func (s *PostService) FollowUnFollow(id uint, jwt *models.JWTDto) *models.SingleApiResponse {
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

// func (s *PostService) ReactionPost(id uint, jwt *models.JWTDto) *models.SingleApiResponse {
// 	if s.db.Where("profile_id =? AND post_id =?", jwt.ID).First(&models.Reaction{}).Error != nil {
// 		s.db.Create(&models.Reaction{ProfileID: jwt.ID, PostID: id})
// 	} else {
// 		s.db.Where("profile_id =? AND post_id =?", jwt.ID).Delete(&models.Reaction{})
// 	}

// 	return models.NewSingleOkSingle()
// }

// func NewPostService(db *gorm.DB) *PostService {
// 	return &PostService{db: db}
// }
