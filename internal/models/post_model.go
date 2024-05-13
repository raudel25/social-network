package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title           string
	ProfileID         uint
	Profile           User `gorm:"foreignKey:ProfileID"`
	PhotoID *uint
	Photo   *Photo `gorm:"foreignKey:PhotoID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	RichText       RichText `gorm:"type:jsonb"`
	RePostID *uint
	RePost   *Post `gorm:"foreignKey:RePostID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
