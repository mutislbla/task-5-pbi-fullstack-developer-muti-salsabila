package app

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	Title    string `json:"title" valid:"required"`
	Caption  string `json:"caption" valid:"required"`
	PhotoUrl string `json:"photo_url" valid:"required"`
	UserID   uint   `json:"userId"`
}
