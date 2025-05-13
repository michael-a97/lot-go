package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (request ChangePasswordRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(
			&request.OldPassword,
			validation.Required.Error("old_password is required"),
		),
		validation.Field(
			&request.NewPassword,
			validation.Required.Error("new_password is required"),
		),
	)
}
