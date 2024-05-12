package models

import (
	"gorm.io/gorm"
)

type Reaction struct {
	gorm.Model
	UserID   uint
	User     User     `gorm:"foreignKey:UserID"`
	PostID   uint
}
