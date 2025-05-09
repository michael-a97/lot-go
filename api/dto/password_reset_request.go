package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type PasswordResetRequest struct {
	PhoneNumber string `json:"phone_number"`
	NewPassword string `json:"new_password"`
}

func (request PasswordResetRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(
			&request.NewPassword,
			validation.Required.Error("new_password is required"),
		),
		validation.Field(
			&request.PhoneNumber,
			validation.Required.Error("phone_number is required"),
		),
		
	)
}
