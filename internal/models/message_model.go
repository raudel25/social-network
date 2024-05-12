package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserID   uint
	User     User     `gorm:"foreignKey:UserID"`
	RichText RichText `gorm:"type:jsonb"`
	PostID   uint
}
