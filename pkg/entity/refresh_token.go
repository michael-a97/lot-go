package entity

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model
	Token   string `json:"token"`
	Revoked bool   `json:"revoked"`
	UserID uint 
}
