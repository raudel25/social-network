package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	Filename string `json:"filename"`
}
