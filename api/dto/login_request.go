package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (request LoginRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(
			&request.Password,
			validation.Required.Error("password is required"),
		),
		validation.Field(
			&request.PhoneNumber,
			validation.Required.Error("phone number is required"),
		),
	)
}
