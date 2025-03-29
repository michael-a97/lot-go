package dto

type UserDto struct {
	PhoneNumber string `json:"phone_number" gorm:"unique"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Role        string `json:"role"`
}
