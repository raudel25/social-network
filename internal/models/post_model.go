package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title     string
	ProfileID uint
	Profile   Profile `gorm:"foreignKey:ProfileID"`
	PhotoID   *uint
	RichText  RichText `gorm:"type:jsonb"`
	RePostID  *uint
	RePost    *Post      `gorm:"foreignKey:RePostID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Messages  []Message  `gorm:"foreignKey:PostID"`
	Reactions []Reaction `gorm:"foreignKey:PostID"`
	Seen []SeenPost `gorm:"foreignKey:PostID"`
}
