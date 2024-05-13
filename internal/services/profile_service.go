package services

import (
	"social-network-api/internal/models"

	"gorm.io/gorm"
)

type ProfileService struct {
	db *gorm.DB
}

func profileToResponseProfile(id uint, profile *models.Profile) *models.ProfileResponse {
	follow := false
	for _, v := range profile.FollowedBy {
		if v.FollowerProfileID == id {
			follow = true
			break
		}

	}
	return &models.ProfileResponse{
		ID:           profile.ID,
		Name:         profile.Name,
		ProfilePhoto: profile.ProfilePhoto,
		BannerPhoto:  profile.BannerPhoto,
		RichText:     profile.RichText,
		Follow:       follow,
		Username:     profile.User.Username,
	}
}

func (s *ProfileService) GetByFollowed(id uint, jwt *models.JWTDto) *models.ApiResponse[[]models.ProfileResponse] {
	var followerProfiles []models.Profile

	s.db.Table("follows").Select("*").
		Joins("join profiles on follows.follower_profile_id = profiles.id").
		Where("follows.followed_profile_id =?", id).Preload("FollowedBy").Preload("User").Preload("ProfilePhoto").Preload("BannerPhoto").
		Find(&followerProfiles)

	var profiles []models.ProfileResponse

	for _, v := range followerProfiles {
		profiles = append(profiles, *profileToResponseProfile(jwt.ID, &v))
	}

	return models.NewOk(&profiles)
}

func (s *ProfileService) GetByFollower(id uint, jwt *models.JWTDto) *models.ApiResponse[[]models.ProfileResponse] {
	var followedProfiles []models.Profile

	s.db.Table("follows").Select("*").
		Joins("join profiles on follows.followed_profile_id = profiles.id").
		Where("follows.follower_profile_id =?", id).Preload("FollowedBy").Preload("User").Preload("ProfilePhoto").Preload("BannerPhoto").
		Find(&followedProfiles)

	var profiles []models.ProfileResponse

	for _, v := range followedProfiles {
		profiles = append(profiles, *profileToResponseProfile(jwt.ID, &v))
	}

	return models.NewOk(&profiles)
}

// func (s *ProfileService) GetByName(name string, jwt *models.JWTDto) *models.ApiResponse[[]models.ProfileResponse] {
// 	var nameProfiles []models.Profile

// 	s.db.Preload("FollowedBy").Preload("User").Preload("ProfilePhoto").Preload("BannerPhoto").Where("name =?", name).Find(&nameProfiles)

// 	var profiles []models.ProfileResponse

// 	for _, v := range nameProfiles {
// 		profiles = append(profiles, *profileToResponseProfile(jwt.ID, &v))
// 	}

// 	return models.NewOk(&profiles)
// }

func (s *ProfileService) GetByID(id uint, jwt *models.JWTDto) *models.ApiResponse[models.ProfileResponse] {
	var profile models.Profile
	if s.db.Preload("FollowedBy").Preload("User").Preload("ProfilePhoto").Preload("BannerPhoto").First(&profile, id).Error != nil {
		return models.NewNotFound[models.ProfileResponse]("Profile not found")
	}

	return models.NewOk(profileToResponseProfile(jwt.ID, &profile))
}

func (s *ProfileService) EditProfile(request *models.ProfileRequest, jwt *models.JWTDto) *models.SingleApiResponse {
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
