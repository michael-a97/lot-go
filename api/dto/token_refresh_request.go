package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (request TokenRefreshRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(
			&request.RefreshToken,
			validation.Required.Error("Refresh token is required"),
		),
	)
}
