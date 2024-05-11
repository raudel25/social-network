package models

import "gorm.io/gorm"
type User struct {
	gorm.Model
	Id       uint `gorm:"primaryKey"`
	Name     string
	UserName string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}
