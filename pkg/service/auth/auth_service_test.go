package service

import (
	dto "lot/api/dto/auth"
	"lot/pkg/entity"
	app_errors "lot/pkg/errors"
	authRepository "lot/pkg/repository/auth"
	userRepository "lot/pkg/repository/user"
	service "lot/pkg/service/sms_token_verifier"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestSignIn(t *testing.T) {
	t.Run("should return account not found error"+
		" when phone number is invalid",
		func(t *testing.T) {
			mockRepository := authRepository.MockAuthRepository{}
			mockUserRepository := userRepository.MockUserRepository{}
			smsTokenVerifier := service.MockSmsTokenVerifier{}
			signInRequest := dto.LoginRequest{
				PhoneNumber: "+251923001100",
				Password:    "password",
			}
			mockUserRepository.On(
				"FindByPhoneNumber",
				signInRequest.PhoneNumber,
			).Return(
				&entity.User{},
				app_errors.ErrAccountNotFound,
			)
			authService := NewAuthService(
				&mockRepository,
				&mockUserRepository,
				&smsTokenVerifier,
			)

			response, err := authService.SignIn(signInRequest)

			assert.Equal(t, app_errors.ErrAccountNotFound, err)
			assert.Nil(t, response)

		})

	t.Run("should return invalid phone number or password error"+
		" when password is incorrect",
		func(t *testing.T) {
			mockRepository := authRepository.MockAuthRepository{}
			mockUserRepository := userRepository.MockUserRepository{}
			smsTokenVerifier := service.MockSmsTokenVerifier{}
			signInRequest := dto.LoginRequest{
				PhoneNumber: "+251923001100",
				Password:    "password",
			}
			invalidPassword := "$2a$14$KbKNc.ZXCZAZ2jv/0KF4/OvmpKDFNDxpMMmeYG8OiSX/9k0EVOEvX"
			savedUser := entity.User{
				PhoneNumber: signInRequest.PhoneNumber,
				Password:    invalidPassword,
				// Password: "$2a$14$KbKNc.ZXCZAZ2jv/0KF4/OvmpKDFNDxpMMmeYG8OiSX/9k0EVOEvW",//Hashed value for the string "password"
				FirstName: "John",
				LastName:  "Doe",
			}
			mockUserRepository.On(
				"FindByPhoneNumber",
				signInRequest.PhoneNumber,
			).Return(&savedUser, nil)
			authService := NewAuthService(
				&mockRepository,
				&mockUserRepository,
				&smsTokenVerifier,
			)

			response, err := authService.SignIn(signInRequest)

			assert.Equal(t, app_errors.ErrInvalidPhoneNumberOrPassword, err)
			assert.Nil(t, response)

		})

	t.Run("should return error when unable to get signed token",
		func(t *testing.T) {
			mockAuthRepository := authRepository.MockAuthRepository{}
			mockUserRepository := userRepository.MockUserRepository{}
			smsTokenVerifier := service.MockSmsTokenVerifier{}
			signInRequest := dto.LoginRequest{
				PhoneNumber: "+251923001100",
				Password:    "password",
			}
			validPassword := "$2a$14$KbKNc.ZXCZAZ2jv/0KF4/OvmpKDFNDxpMMmeYG8OiSX/9k0EVOEvW"
			savedUser := entity.User{
				PhoneNumber: signInRequest.PhoneNumber,
				Password:    validPassword,
				FirstName:   "John",
				LastName:    "Doe",
				Role: entity.Role{
					Model: gorm.Model{ID: 1},
					Name:  entity.RoleAttendant,
				},
			}
			mockUserRepository.On(
				"FindByPhoneNumber",
				signInRequest.PhoneNumber,
			).Return(&savedUser, nil)
			authService := NewAuthService(
				&mockAuthRepository,
				&mockUserRepository,
				&smsTokenVerifier,
			)
			mockAuthRepository.On("RevokeAllRefreshTokensForUser", savedUser.ID).Return(nil)
			mockAuthRepository.On("SaveRefreshToken", mock.Anything, savedUser.ID).Return(nil)

			_, err := authService.SignIn(signInRequest)

			assert.Error(t, err)
		},
	)

	t.Run("should return authentication response when successful",
		func(t *testing.T) {
			os.Setenv("secret", "1234")
			mockAuthRepository := authRepository.MockAuthRepository{}
			mockUserRepository := userRepository.MockUserRepository{}
			smsTokenVerifier := service.MockSmsTokenVerifier{}
			signInRequest := dto.LoginRequest{
				PhoneNumber: "+251923001100",
				Password:    "password",
			}
			validPassword := "$2a$14$KbKNc.ZXCZAZ2jv/0KF4/OvmpKDFNDxpMMmeYG8OiSX/9k0EVOEvW"
			savedUser := entity.User{
				PhoneNumber: signInRequest.PhoneNumber,
				Password:    validPassword,
				FirstName:   "John",
				LastName:    "Doe",
				Role: entity.Role{
					Model: gorm.Model{ID: 1},
					Name:  entity.RoleAttendant,
				},
			}
			mockUserRepository.On(
				"FindByPhoneNumber",
				signInRequest.PhoneNumber,
			).Return(&savedUser, nil)
			authService := NewAuthService(
				&mockAuthRepository,
				&mockUserRepository,
				&smsTokenVerifier,
			)
			mockAuthRepository.On("RevokeAllRefreshTokensForUser", savedUser.ID).Return(nil)
			mockAuthRepository.On("SaveRefreshToken", mock.Anything, savedUser.ID).Return(nil)

			response, err := authService.SignIn(signInRequest)

			assert.Nil(t, err)
			assert.IsType(t, dto.AuthenticationResponse{}, *response)

		})

}
