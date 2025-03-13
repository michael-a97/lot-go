package dto

import (
	"fmt"
	"lot/pkg/entity"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SignUpRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

func (request SignUpRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(
			&request.FirstName,
			validation.Required.Error("First name is required"),
		),
		validation.Field(
			&request.LastName,
			validation.Required.Error("Last name is required"),
		),
		validation.Field(
			&request.Password,
			validation.Required.Error("Password is required"),
		),
		validation.Field(
			&request.PhoneNumber,
			validation.Required.Error("Phone number is required"),
		),
		validation.Field(
			&request.Role,
			validation.Required.Error("Role is required"),
			validation.In(entity.RoleAttendant, entity.RoleUser).Error(
				fmt.Sprintf(
					"role can be one of the following [\"%s\" \"%s\"]",
					entity.RoleUser, entity.RoleAttendant,
				),
			),
		),
	)
}
