package service

import (
	"errors"
	"log"
	"lot/api/dto"
	"lot/internal/entity"
	"lot/internal/errors"
	repositoryMocks "lot/internal/repository/mocks"
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
			mockRepository := repositoryMocks.MockAuthRepository{}
			mockUserRepository := repositoryMocks.MockUserRepository{}
			smsTokenVerifier := MockSmsTokenVerifier{}
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
			mockRepository := repositoryMocks.MockAuthRepository{}
			mockUserRepository := repositoryMocks.MockUserRepository{}
			smsTokenVerifier := MockSmsTokenVerifier{}
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
			mockAuthRepository := repositoryMocks.MockAuthRepository{}
			mockUserRepository := repositoryMocks.MockUserRepository{}
			smsTokenVerifier := MockSmsTokenVerifier{}
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
			if err := os.Setenv("secret", "1234"); err != nil {
				log.Fatal(err.Error())
			}
			mockAuthRepository := repositoryMocks.MockAuthRepository{}
			mockUserRepository := repositoryMocks.MockUserRepository{}
			smsTokenVerifier := MockSmsTokenVerifier{}
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

func TestResetPassword(t *testing.T) {
	t.Run("should return error when the phone number in the request is not found",
		func(t *testing.T) {
			mockAuthRepository := repositoryMocks.MockAuthRepository{}
			mockUserRepository := repositoryMocks.MockUserRepository{}
			smsTokenVerifier := MockSmsTokenVerifier{}

			request := dto.PasswordResetRequest{
				PhoneNumber: "+251923001100",
				NewPassword: "my-new-passw0rd",
			}

			authService := NewAuthService(
				&mockAuthRepository,
				&mockUserRepository,
				&smsTokenVerifier,
			)

			var nilUserPointer *entity.User

			mockUserRepository.On(
				"FindByPhoneNumber",
				request.PhoneNumber,
			).Return(nilUserPointer, app_errors.ErrRecordNotFound)

			err := authService.ResetPassword(request)

			assert.Error(t, err)

		},
	)

	t.Run("should return nill when successful",
		func(t *testing.T) {
			mockAuthRepository := repositoryMocks.MockAuthRepository{}
			mockUserRepository := repositoryMocks.MockUserRepository{}
			smsTokenVerifier := MockSmsTokenVerifier{}
			request := dto.PasswordResetRequest{
				PhoneNumber: "+251923001100",
				NewPassword: "my-new-passw0rd",
			}
			authService := NewAuthService(
				&mockAuthRepository,
				&mockUserRepository,
				&smsTokenVerifier,
			)
			user := entity.User{FirstName: "John", LastName: "Doe"}
			mockUserRepository.On("FindByPhoneNumber", request.PhoneNumber).Return(&user, nil)
			mockUserRepository.On(
				"Save",
				mock.Anything,
			).Return(&user, nil)

			err := authService.ResetPassword(request)

			assert.Nil(t, err)
		},
	)

}

func TestChangePassword(t *testing.T) {
	t.Run("should return invalid password error ", func(t *testing.T) {
		mockAuthRepository := repositoryMocks.MockAuthRepository{}
		mockUserRepository := repositoryMocks.MockUserRepository{}
		smsTokenVerifier := MockSmsTokenVerifier{}

		request := dto.ChangePasswordRequest{
			OldPassword: "2222",
			NewPassword: "1234",
		}
		user := entity.User{
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "+251923001100",
			Password:    "$2a$14$A/zW6mNZEnK8s1zIKzSl0uo1yNuQtizeeMQA.HrwxxpyR37AzwV3x",
		}

		authService := NewAuthService(
			&mockAuthRepository,
			&mockUserRepository,
			&smsTokenVerifier,
		)

		err := authService.ChangePassword(request, user)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), app_errors.ErrInvalidPassword.Error())
	})

	t.Run("should return error if there's an error when updating user", func(t *testing.T) {
		mockAuthRepository := repositoryMocks.MockAuthRepository{}
		mockUserRepository := repositoryMocks.MockUserRepository{}
		smsTokenVerifier := MockSmsTokenVerifier{}

		request := dto.ChangePasswordRequest{
			OldPassword: "2222",
			NewPassword: "1234",
		}
		user := entity.User{
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "+251923001100",
			Password:    "$2a$14$A/zW6mNZEnK8s1zIKzSl0uo1yNuQtizeeMQA.HrwxxpyR37AzwV3q",
		}

		mockUserRepository.On("Update", mock.Anything).Return(&entity.User{}, errors.New("some update db error"))

		authService := NewAuthService(
			&mockAuthRepository,
			&mockUserRepository,
			&smsTokenVerifier,
		)

		err := authService.ChangePassword(request, user)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errors.New("some update db error").Error())
	})

}
