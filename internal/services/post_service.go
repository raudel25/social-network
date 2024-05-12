package services

import (
	"social-network-api/internal/models"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func (s *PostService) NewPost(request *models.PostRequest) *models.SingleApiResponse {
	return models.NewSingleOkSingle()
}

func (s *PostService) MessagePost(request *models.MessageRequest, id uint, jwt *models.JWTDto) *models.SingleApiResponse {
	return models.NewSingleOkSingle()
}

func (s *PostService) ReactionPost(id uint, jwt *models.JWTDto) *models.SingleApiResponse {
	return models.NewSingleOkSingle()
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}
