package models

import (
	"gorm.io/gorm"
)

type SeenPost struct {
	gorm.Model
	ProfileID uint
	Profile   User `gorm:"foreignKey:ProfileID"`
	PostID    uint
}
