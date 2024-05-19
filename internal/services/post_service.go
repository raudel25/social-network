package services

import (
	"fmt"
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
		RichText:      post.RichText,
		Reaction:      reaction,
		CantReactions: len(post.Reactions),
		CantMessages:  len(post.Messages),
		RePost:        rePost,
		Date:          post.CreatedAt,
	}
}

func (s *PostService) GetByRecommendationPost(pagination *pkg.Pagination[models.PostResponse], jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination[models.PostResponse]] {
	var recommendationPosts []models.Post

	query := fmt.Sprintf(`
	SELECT *
	FROM posts
	WHERE 
	EXISTS (SELECT 1 FROM follows WHERE follower_profile_id=%d AND followed_profile_id = profile_id)
	AND 
	NOT EXISTS (SELECT * FROM seen_posts WHERE profile_id=%d AND created_at < now() - interval '24 hours')`,
		jwt.ID, jwt.ID)

	pagination.CountRaw(s.db, query)
	s.db.Raw(query).Scopes(pagination.Paginate).Scan(&recommendationPosts)

	var posts []models.PostResponse

	for _, v := range recommendationPosts {
		s.db.Preload("Reactions").Preload("Messages").Preload("Profile").
			Preload("Profile.User").Preload("Profile.FollowedBy").Find(&v, v.ID)
		posts = append(posts, *postToResponsePost(jwt.ID, &v))
	}

	pagination.Rows = posts

	return pkg.NewOk(pagination)
}

func (s *PostService) GetPostsByUser(pagination *pkg.Pagination[models.PostResponse], id uint, jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination[models.PostResponse]] {
	var posts []models.Post
	pagination.Count(s.db.Where("profile_id =?", id).Model(&models.Post{}))

	s.db.Where("profile_id =?", id).Scopes(pagination.Paginate).
		Preload("Reactions").Preload("Messages").Preload("Profile").
		Preload("Profile.User").Preload("Profile.FollowedBy").Find(&posts)

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
