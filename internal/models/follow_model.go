package models

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	FollowerProfileID uint 
	FollowedProfileID uint
}
