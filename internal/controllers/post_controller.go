package controllers

import (
	"social-network-api/internal/models"
	"social-network-api/internal/services"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	jwtService  *services.JWTService
	postService *services.PostService
}

func (s *PostController) GetPost(c *gin.Context) {
	checkAuthorized := CheckAuthorized(c, s.jwtService)

	if !checkAuthorized.Ok() {
		checkAuthorized.Response(c)
		return
	}

	checkPagination := CheckPagination[models.PostResponse](c)
	if !checkPagination.Ok() {
		checkPagination.Response(c)
		return
	}

	s.postService.GetByRecommendationPost(checkPagination.Data, checkAuthorized.Data).Response(c)
}

func (s *PostController) GetPostsByUser(c *gin.Context) {
	checkAuthorized := CheckAuthorized(c, s.jwtService)

	if !checkAuthorized.Ok() {
		checkAuthorized.Response(c)
		return
	}

	idRequest := CheckId(c)
	if !idRequest.Ok() {
		idRequest.Response(c)
		return
	}

	checkPagination := CheckPagination[models.PostResponse](c)
	if !checkPagination.Ok() {
		checkPagination.Response(c)
		return
	}

	s.postService.GetPostsByUser(checkPagination.Data, *idRequest.Data, checkAuthorized.Data).Response(c)
}

func (s *PostController) NewPost(c *gin.Context) {
	checkAuthorized := CheckAuthorized(c, s.jwtService)

	if !checkAuthorized.Ok() {
		checkAuthorized.Response(c)
		return
	}

	var request models.PostRequest

	checkRequest := CheckRequest(c, &request)
	if !checkRequest.Ok() {
		checkRequest.Response(c)
		return
	}

	s.postService.NewPost(&request, checkAuthorized.Data).Response(c)
}

func (s *PostController) MessagePost(c *gin.Context) {
	checkAuthorized := CheckAuthorized(c, s.jwtService)

	if !checkAuthorized.Ok() {
		checkAuthorized.Response(c)
		return
	}

	var request models.MessageRequest

	checkRequest := CheckRequest(c, &request)
	if !checkRequest.Ok() {
		checkRequest.Response(c)
		return
	}

	idRequest := CheckId(c)
	if !idRequest.Ok() {
		idRequest.Response(c)
		return
	}

	s.postService.MessagePost(&request, *idRequest.Data, checkAuthorized.Data).Response(c)
}

func (s *PostController) ReactionPost(c *gin.Context) {
	checkAuthorized := CheckAuthorized(c, s.jwtService)

	if !checkAuthorized.Ok() {
		checkAuthorized.Response(c)
		return
	}

	idRequest := CheckId(c)
	if !idRequest.Ok() {
		idRequest.Response(c)
		return
	}

	s.postService.ReactionPost(*idRequest.Data, checkAuthorized.Data).Response(c)

}

func NewPostController(postService *services.PostService, jwtService *services.JWTService) *PostController {
	return &PostController{jwtService: jwtService, postService: postService}
}
