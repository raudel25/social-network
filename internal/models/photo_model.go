package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	Name string
	Src  string
}
