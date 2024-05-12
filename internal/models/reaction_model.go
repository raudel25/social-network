package models

import (
	"gorm.io/gorm"
)

type Reaction struct {
	gorm.Model
	ProfileID uint
	Profile   User `gorm:"foreignKey:ProfileID"`
	PostID    uint
}
