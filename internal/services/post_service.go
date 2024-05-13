package services

import (
	"social-network-api/internal/models"
	"social-network-api/internal/pkg"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func (s *PostService) NewPost(request *models.PostRequest, jwt *models.JWTDto) *pkg.SingleApiResponse {
	if request.PhotoID != nil && s.db.First(&models.Photo{}, request.PhotoID).Error != nil {
		return pkg.NewSingleNotFound("Photo not found")
	}

	if request.RePostID != nil && s.db.First(&models.Post{}, request.RePostID).Error != nil {
		return pkg.NewSingleNotFound("Post not found")
	}

	s.db.Create(&models.Post{Title: request.Title, ProfileID: jwt.ID, PhotoID: request.PhotoID, RePostID: request.RePostID, RichText: request.RichText})

	return pkg.NewSingleOkSingle()
}

func (s *PostService) MessagePost(request *models.MessageRequest, id uint, jwt *models.JWTDto) *pkg.SingleApiResponse {
	if s.db.First(&models.Post{}, id).Error != nil {
		return pkg.NewSingleNotFound("Post not found")
	}

	s.db.Create(&models.Message{ProfileID: jwt.ID, RichText: request.RichText, PostID: id})

	return pkg.NewSingleOkSingle()
}

func (s *PostService) ReactionPost(id uint, jwt *models.JWTDto) *pkg.SingleApiResponse {
	if s.db.Where("profile_id =? AND post_id =?", jwt.ID).First(&models.Reaction{}).Error != nil {
		s.db.Create(&models.Reaction{ProfileID: jwt.ID, PostID: id})
	} else {
		s.db.Where("profile_id =? AND post_id =?", jwt.ID).Delete(&models.Reaction{})
	}

	return pkg.NewSingleOkSingle()
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}
