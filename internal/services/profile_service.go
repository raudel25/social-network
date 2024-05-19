package services

import (
	"fmt"
	"social-network-api/internal/models"
	"social-network-api/internal/pkg"

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
		ID:       profile.ID,
		Name:     profile.Name,
		RichText: profile.RichText,
		Follow:   follow,
		Username: profile.User.Username,
	}
}

func (s *ProfileService) GetByFollowed(pagination *pkg.Pagination[models.ProfileResponse], username string, jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination[models.ProfileResponse]] {
	var followerProfiles []models.Profile

	response := s.GetByUsername(username, jwt)
	if !response.Ok() {
		return pkg.NewNotFound[pkg.Pagination[models.ProfileResponse]]("Profile not found")
	}

	id := response.Data.ID

	pagination.Count(s.db.Table("follows").Select("*").
		Joins("join profiles on follows.follower_profile_id = profiles.id").
		Where("follows.followed_profile_id =?", id))

	s.db.Table("follows").Select("*").
		Joins("join profiles on follows.follower_profile_id = profiles.id").
		Where("follows.followed_profile_id =?", id).Scopes(pagination.Paginate).
		Preload("FollowedBy").Preload("User").
		Find(&followerProfiles)

	var profiles []models.ProfileResponse

	for _, v := range followerProfiles {
		profiles = append(profiles, *profileToResponseProfile(jwt.ID, &v))
	}

	pagination.Rows = profiles

	return pkg.NewOk(pagination)
}

func (s *ProfileService) GetByFollower(pagination *pkg.Pagination[models.ProfileResponse], username string, jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination[models.ProfileResponse]] {
	var followedProfiles []models.Profile

	response := s.GetByUsername(username, jwt)
	if !response.Ok() {
		return pkg.NewNotFound[pkg.Pagination[models.ProfileResponse]]("Profile not found")
	}

	id := response.Data.ID

	pagination.Count(s.db.Table("follows").Select("*").
		Joins("join profiles on follows.followed_profile_id = profiles.id").
		Where("follows.follower_profile_id =?", id))

	s.db.Table("follows").Select("*").
		Joins("join profiles on follows.followed_profile_id = profiles.id").
		Where("follows.follower_profile_id =?", id).
		Scopes(pagination.Paginate).
		Preload("FollowedBy").Preload("User").
		Find(&followedProfiles)

	var profiles []models.ProfileResponse

	for _, v := range followedProfiles {
		profiles = append(profiles, *profileToResponseProfile(jwt.ID, &v))
	}

	pagination.Rows = profiles

	return pkg.NewOk(pagination)
}

func (s *ProfileService) GetByRecommendationProfile(pagination *pkg.Pagination[models.ProfileResponse], jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination[models.ProfileResponse]] {
	var recommendationProfiles []models.Profile

	query := fmt.Sprintf(`
	SELECT *
	FROM (
		SELECT f.followed_profile_id AS id  
		FROM follows as f
		WHERE 
		EXISTS (SELECT 1 FROM follows WHERE follower_profile_id=%d AND followed_profile_id = f.follower_profile_id)
		AND
		NOT EXISTS (SELECT 1 FROM follows WHERE follower_profile_id=%d AND followed_profile_id = f.followed_profile_id)
		AND
		f.followed_profile_id <> %d
		ORDER BY f.followed_profile_id DESC
	) as f
	JOIN profiles ON f.id=profiles.id`, jwt.ID, jwt.ID, jwt.ID)

	pagination.CountRaw(s.db, query)
	s.db.Raw(query).Scopes(pagination.Paginate).Scan(&recommendationProfiles)

	var profiles []models.ProfileResponse

	for _, v := range recommendationProfiles {
		s.db.Preload("User").Preload("FollowedBy").Find(&v, v.ID)
		profiles = append(profiles, *profileToResponseProfile(jwt.ID, &v))
	}

	pagination.Rows = profiles

	return pkg.NewOk(pagination)
}

func (s *ProfileService) GetByUsername(username string, jwt *models.JWTDto) *pkg.ApiResponse[models.ProfileResponse] {
	var profile models.Profile
	if s.db.Preload("FollowedBy").Preload("User").Where("username =?", username).Joins("JOIN users ON profiles.user_id = users.id").First(&profile).Error != nil {
		return pkg.NewNotFound[models.ProfileResponse]("Profile not found")
	}

	return pkg.NewOk(profileToResponseProfile(jwt.ID, &profile))
}

func (s *ProfileService) EditProfile(request *models.ProfileRequest, jwt *models.JWTDto) *pkg.SingleApiResponse {
	if request.ProfilePhotoID != nil && s.db.First(&models.Photo{}, request.ProfilePhotoID).Error != nil {
		return pkg.NewSingleNotFound("Profile photo not found")
	}

	if request.BannerPhotoID != nil && s.db.First(&models.Photo{}, request.BannerPhotoID).Error != nil {
		return pkg.NewSingleNotFound("Banner photo not found")
	}

	var profile models.Profile
	if s.db.Find(&profile, jwt.ID).Error != nil {
		return pkg.NewSingleNotFound("Profile not found")
	}

	profile.Name = request.Name
	profile.ProfilePhotoID = request.ProfilePhotoID
	profile.BannerPhotoID = request.BannerPhotoID
	profile.RichText = request.RichText

	s.db.Where("id =?", jwt.ID).Updates(&profile)

	return pkg.NewSingleOkSingle()
}

func (s *ProfileService) FollowUnFollow(id uint, jwt *models.JWTDto) *pkg.SingleApiResponse {
	if s.db.First(&models.Profile{}, id).Error != nil {
		return pkg.NewSingleNotFound("Profile not found")
	}

	if s.db.Where("follower_id =? AND followed_id =?", jwt.ID, id).First(&models.Follow{}).Error != nil {
		s.db.Create(&models.Follow{FollowerProfileID: jwt.ID, FollowedProfileID: id})
	} else {
		s.db.Where("follower_id =? AND followed_id =?", jwt.ID, id).Delete(&models.Follow{})
	}

	return pkg.NewSingleOkSingle()
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{db: db}
}
