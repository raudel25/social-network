package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}
