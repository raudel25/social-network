package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"social-network-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PhotoController struct {
	photoService *services.PhotoService
	jwtService   *services.JWTService
}

func NewPhotoController(photoService *services.PhotoService, jwtService *services.JWTService) *PhotoController {
	return &PhotoController{photoService: photoService, jwtService: jwtService}
}

func (s *PhotoController) GetPhoto(c *gin.Context) {
	id := CheckId(c)
	if !id.Ok() {
		id.Response(c)
		return
	}

	response := s.photoService.GetPhoto(*id.Data)
	if !response.Ok() {
		response.Response(c)
		return
	}

	filePath := response.Data.Filename

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"ok": false, "message": "Photo not found"})
		return
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "message": err.Error()})
		return
	}

	c.Data(http.StatusOK, http.DetectContentType(fileBytes), fileBytes)
}

func (s *PhotoController) UploadPhoto(c *gin.Context) {
	checkAuthorized := CheckAuthorized(c, s.jwtService)
	if !checkAuthorized.Ok() {
		checkAuthorized.Response(c)
		return
	}

	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "message": err.Error()})
		return
	}

	newFileName := uuid.New().String()
	savePath := filepath.Join("uploads", newFileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "message": err.Error()})
		return
	}

	s.photoService.UploadPhoto(savePath).Response(c)
}
