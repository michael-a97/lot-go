package entity

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model
	Token   string `json:"token" gorm:"required"`
	Revoked bool   `json:"revoked" gorm:"required"`
	UserID uint 
}
