package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title           string
	UserID         uint
	User           User `gorm:"foreignKey:UserID"`
	PhotoID *uint
	Photo   *Photo `gorm:"foreignKey:PhotoID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	RichText       RichText `gorm:"type:jsonb"`
	RePostID *uint
	RePost   *Post `gorm:"foreignKey:RePostID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
