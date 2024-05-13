package controllers

import (
	"social-network-api/internal/models"
	"social-network-api/internal/services"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	jwtService     *services.JWTService
	profileService *services.ProfileService
}

func (s *ProfileController) GetByRecommendation(c *gin.Context) {
	checkAuthorized := CheckAuthorized(c, s.jwtService)

	if !checkAuthorized.Ok() {
		checkAuthorized.Response(c)
		return
	}

	s.profileService.GetByRecommendation(checkAuthorized.Data).Response(c)
}

func (s *ProfileController) GetByFollowed(c *gin.Context) {
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

	s.profileService.GetByFollowed(*idRequest.Data, checkAuthorized.Data).Response(c)
}

func (s *ProfileController) GetByFollower(c *gin.Context) {
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

	s.profileService.GetByFollower(*idRequest.Data, checkAuthorized.Data).Response(c)
}

func (s *ProfileController) GetByID(c *gin.Context) {
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

	s.profileService.GetByID(*idRequest.Data, checkAuthorized.Data).Response(c)
}

func (s *ProfileController) EditProfile(c *gin.Context) {
	checkAuthorized := CheckAuthorized(c, s.jwtService)

	if !checkAuthorized.Ok() {
		checkAuthorized.Response(c)
		return
	}

	var request models.ProfileRequest

	checkRequest := CheckRequest(c, &request)
	if !checkRequest.Ok() {
		checkRequest.Response(c)
		return
	}

	s.profileService.EditProfile(&request, checkAuthorized.Data).Response(c)
}

func (s *ProfileController) FollowUnFollow(c *gin.Context) {
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

	s.profileService.FollowUnFollow(*idRequest.Data, checkAuthorized.Data).Response(c)
}

func NewProfileController(profileService *services.ProfileService, jwtService *services.JWTService) *ProfileController {
	return &ProfileController{jwtService: jwtService, profileService: profileService}
}
