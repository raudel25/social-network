package models

type ProfileRequest struct {
	Name           string    `json:"name" binding:"required"`
	ProfilePhotoID *uint     `json:"profilePhotoId"`
	BannerPhotoID  *uint     `json:"bannerPhotoId"`
	RichText       *RichText `json:"rich_text" binding:"required"`
}

type ProfileResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	ProfilePhotoID *uint     `json:"profilePhotoId"`
	BannerPhotoID  *uint     `json:"bannerPhotoId"`
	RichText     *RichText `json:"richText"`
	Follow       bool      `json:"follow"`
	Username     string    `json:"username"`
}
