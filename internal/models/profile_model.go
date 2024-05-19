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
	BannerPhotoID  *uint
	RichText       *RichText `gorm:"type:jsonb"`
	Follows        []Follow  `gorm:"foreignKey:FollowerProfileID"`
	FollowedBy     []Follow  `gorm:"foreignKey:FollowedProfileID"`
}
