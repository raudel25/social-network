package models

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	Name           string
	UserID         uint
	User           *User `gorm:"foreignKey:UserID"`
	ProfilePhotoID *uint
	ProfilePhoto   *Photo `gorm:"foreignKey:ProfilePhotoID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	BannerPhotoID  *uint
	BannerPhoto    *Photo   `gorm:"foreignKey:BannerPhotoID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	RichText       RichText `gorm:"type:jsonb"`
}
