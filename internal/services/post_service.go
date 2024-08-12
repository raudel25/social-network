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
		if v.ProfileID == id {
			reaction = true
		}
	}

	var rePost *models.PostResponse
	if post.RePost != nil {
		rePost = postToResponsePost(id, post.RePost)
	}

	var messages []models.MessageResponse
	for _, m := range post.Messages {
		messages = append(messages, models.MessageResponse{ID: m.ID, RichText: m.RichText, Profile: *profileToResponseProfile(id, &m.Profile), Date: m.CreatedAt})
	}

	return &models.PostResponse{
		ID:            post.ID,
		Profile:       *profileToResponseProfile(id, &post.Profile),
		RichText:      post.RichText,
		Reaction:      reaction,
		CantReactions: post.CantReactions,
		CantMessages:  post.CantMessages,
		CantRePosts:   post.CantRePosts,
		RePost:        rePost,
		PhotoID:       post.PhotoID,
		Messages:      messages,
		Date:          post.CreatedAt,
	}
}

func (s *PostService) getByRecommendationPost(pagination *pkg.Pagination[models.PostResponse], jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination[models.PostResponse]] {
	var recommendationPosts []models.Post

	query := fmt.Sprintf(`
	SELECT *
	FROM posts
	WHERE 
	EXISTS (SELECT 1 FROM follows WHERE follower_profile_id=%d AND followed_profile_id = profile_id)
	AND 
	NOT EXISTS (SELECT * FROM seen_posts WHERE seen_posts.profile_id=%d AND seen_posts.post_id=posts.id AND created_at < now() - interval '24 hours')`,
		jwt.ID, jwt.ID)

	pagination.CountRaw(s.db, query)
	s.db.Raw(pagination.PaginateRaw(query)).Scan(&recommendationPosts)

	var posts []models.PostResponse

	for _, v := range recommendationPosts {
		s.db.Preload("Reactions").Preload("Profile").
			Preload("Profile.User").Preload("Profile.FollowedBy").Find(&v, v.ID)
		posts = append(posts, *postToResponsePost(jwt.ID, &v))
	}

	pagination.Rows = posts

	return pkg.NewOk(pagination)
}

func (s *PostService) GetByRecommendationPost(pagination *pkg.Pagination[models.PostResponse], jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination[models.PostResponse]] {
	currentPage := s.getByRecommendationPost(pagination, jwt)

	if !currentPage.Ok() || currentPage.Data.Page == 1 {
		return currentPage
	}

	previewPage := s.getByRecommendationPost(&pkg.Pagination[models.PostResponse]{Page: currentPage.Data.Page - 1, Limit: currentPage.Data.Limit}, jwt)

	if previewPage.Ok() {
		for _, v := range previewPage.Data.Rows {
			if s.db.Where("profile_id =? AND post_id =?", jwt.ID, v.ID).First(&models.SeenPost{}).Error == nil {
				continue
			}

			s.db.Create(&models.SeenPost{ProfileID: jwt.ID, PostID: v.ID})
		}
	}

	return currentPage
}

func (s *PostService) GetPostByID(id uint, jwt *models.JWTDto) *pkg.ApiResponse[models.PostResponse] {
	var post models.Post

	if s.db.Preload("Reactions").Preload("Messages").Preload("Messages.Profile").Preload("Messages.Profile.User").Preload("Profile").
		Preload("Profile.User").Preload("Profile.FollowedBy").First(&post, id).Error != nil {
		return pkg.NewNotFound[models.PostResponse]("Not found post")
	}

	return pkg.NewOk(postToResponsePost(jwt.ID, &post))
}

func (s *PostService) GetPostsByUser(pagination *pkg.Pagination[models.PostResponse], username string, jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination[models.PostResponse]] {
	var profile models.Profile
	if s.db.Preload("User").Where("username =?", username).Joins("JOIN users ON profiles.user_id = users.id").First(&profile).Error != nil {
		return pkg.NewNotFound[pkg.Pagination[models.PostResponse]]("Profile not found")
	}

	id := profile.ID

	var posts []models.Post
	pagination.Count(s.db.Where("profile_id =?", id).Model(&models.Post{}))

	s.db.Where("profile_id =?", id).Scopes(pagination.Paginate).
		Preload("Reactions").
		Preload("Profile").Preload("Profile.User").Preload("Profile.FollowedBy").
		Preload("RePost").Preload("RePost.Reactions").
		Preload("RePost.Profile").Preload("RePost.Profile.User").Preload("RePost.Profile.FollowedBy").Find(&posts)

	var response []models.PostResponse

	for _, v := range posts {
		response = append(response, *postToResponsePost(jwt.ID, &v))
	}

	pagination.Rows = response

	return pkg.NewOk(pagination)
}

func (s *PostService) GetPostsByRePostId(pagination *pkg.Pagination[models.PostResponse], id uint, jwt *models.JWTDto) *pkg.ApiResponse[pkg.Pagination[models.PostResponse]] {
	var posts []models.Post
	pagination.Count(s.db.Where("re_post_id =?", id).Model(&models.Post{}))

	s.db.Where("re_post_id =?", id).Scopes(pagination.Paginate).
		Preload("Reactions").
		Preload("Profile").Preload("Profile.User").Preload("Profile.FollowedBy").
		Preload("RePost").Preload("RePost.Reactions").
		Preload("RePost.Profile").Preload("RePost.Profile.User").Preload("RePost.Profile.FollowedBy").Find(&posts)

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

	if request.RePostID != nil {
		var post models.Post
		if s.db.First(&post, request.RePostID).Error != nil {
			return pkg.NewSingleNotFound("Post not found")
		}

		post.CantRePosts++
		s.db.Where("id =?", request.RePostID).Updates(&post)

	}

	s.db.Create(&models.Post{ProfileID: jwt.ID, PhotoID: request.PhotoID, RePostID: request.RePostID, RichText: request.RichText})

	return pkg.NewSingleOkSingle()
}

func (s *PostService) MessagePost(request *models.MessageRequest, id uint, jwt *models.JWTDto) *pkg.SingleApiResponse {
	var post models.Post
	if s.db.First(&post, id).Error != nil {
		return pkg.NewSingleNotFound("Post not found")
	}

	post.CantMessages++
	s.db.Create(&models.Message{ProfileID: jwt.ID, RichText: request.RichText, PostID: id})
	s.db.Where("id =?", id).Updates(&post)

	return pkg.NewSingleOkSingle()
}

func (s *PostService) ReactionPost(id uint, jwt *models.JWTDto) *pkg.SingleApiResponse {
	var post models.Post
	if s.db.First(&post, id).Error != nil {
		return pkg.NewSingleNotFound("Not fount post")
	}

	if s.db.Where("profile_id =? AND post_id =?", jwt.ID, id).First(&models.Reaction{}).Error != nil {
		s.db.Create(&models.Reaction{ProfileID: jwt.ID, PostID: id})
		post.CantReactions++
	} else {
		s.db.Where("profile_id =? AND post_id =?", jwt.ID, id).Delete(&models.Reaction{})
		post.CantReactions--
	}

	s.db.Where("id =?", id).Updates(&post)

	return pkg.NewSingleOkSingle()
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}
