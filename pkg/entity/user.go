package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PhoneNumber   string `json:"phone_number" gorm:"unique"`
	Password      string `json:"password"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	RoleID        uint
	Role          Role           `gorm:"foreignKey:RoleID"`
	RefreshTokens []RefreshToken `json:"refresh_token"`
}

type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
}

const (
	RoleAttendant = "attendant"
	RoleUser      = "user"
	RoleAdmin     = "admin"
)
