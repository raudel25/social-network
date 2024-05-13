package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ProfileID   uint
	Profile     User     `gorm:"foreignKey:ProfileID"`
	RichText RichText `gorm:"type:jsonb"`
	PostID   uint
}
