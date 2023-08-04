package model

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserId  uint `json:"user_id"`
	VideoId uint `json:"video_id"`
	State   uint
}
