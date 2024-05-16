package models

type ProfileRequest struct {
	Name           string    `json:"name" binding:"required"`
	ProfilePhotoID *uint     `json:"profile_photo_id"`
	BannerPhotoID  *uint     `json:"banner_photo_id"`
	RichText       *RichText `json:"rich_text" binding:"required"`
}

type ProfileResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	ProfilePhoto *Photo    `json:"profile_photo"`
	BannerPhoto  *Photo    `json:"banner_photo"`
	RichText     *RichText `json:"rich_text"`
	Follow       bool      `json:"follow"`
	Username     string    `json:"username"`
}
