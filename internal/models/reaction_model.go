package models

import (
	"gorm.io/gorm"
)

type Reaction struct {
	gorm.Model
	ProfileID uint
	Profile   Profile `gorm:"foreignKey:ProfileID"`
	PostID    uint
}
