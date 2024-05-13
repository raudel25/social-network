package models

type ProfileDto struct {
	Name           string   `json:"name" binding:"required"`
	ProfilePhotoID *uint    `json:"profile_photo_id"`
	BannerPhotoID  *uint    `json:"banner_photo_id"`
	RichText       RichText `json:"rich_text" binding:"required"`
}
