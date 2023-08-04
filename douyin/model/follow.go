package model

import "gorm.io/gorm"

type Following struct {
	gorm.Model
	HostId  uint
	GuestId uint
}

type Followers struct {
	gorm.Model
	HostId  uint
	GuestId uint
}
