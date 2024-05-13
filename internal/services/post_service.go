package services

import (
	"social-network-api/internal/models"
	"social-network-api/internal/pkg"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func postToResponsePost(id uint, post *models.Post) *models.PostResponse {
	reaction := false

	for _, v := range post.Reactions {
		if v.Profile.ID == id {
			reaction = true
		}
	}

	var rePost *models.PostResponse
	if post.RePost != nil {
		rePost = postToResponsePost(id, post.RePost)
	}

	return &models.PostResponse{
		Title:         post.Title,
		Profile:       *profileToResponseProfile(id, &post.Profile),
		Photo:         post.Photo,
		RichText:      post.RichText,
		Reaction:      reaction,
		CantReactions: len(post.Reactions),
		CantMessages:  len(post.Messages),
		RePost:        rePost,
		Date:          post.CreatedAt,
	}
}

func (s *PostService) GetPostsByUser(pagination *pkg.Pagination, id uint, jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination] {
	var posts []models.Post
	pagination.Count(s.db.Where("profile_id =?", id).Model(&models.Post{}))

	s.db.Where("profile_id =?", id).Scopes(pagination.Paginate).
		Preload("Reactions").Preload("Messages").Preload("Profile").Preload("Profile.User").Preload("Profile.FollowedBy").Find(&posts)

	var response []models.PostResponse

	for _, v := range posts {
		response = append(response, *postToResponsePost(jwt.ID, &v))
	}

	pagination.Rows = response

	return pkg.NewOk(pagination)
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
