package entity

import "gorm.io/gorm"

type User struct{
	gorm.Model
	PhoneNumber string `json:"phone_number"`
	Password string `json:"password"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}