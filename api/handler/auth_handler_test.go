package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"lot/api/dto"
	"lot/internal/entity"
	"lot/internal/errors"
	"lot/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestResetPasswordHandler(t *testing.T) {
	path := "/reset-password"
	t.Run("should return 400 when request body is not valid ", func(t *testing.T) {
		app := fiber.New()
		mockAuthService := new(mocks.MockAuthService)
		requestBody := dto.PasswordResetRequest{
			PhoneNumber: "+251923001100",
			NewPassword: "",
		}

		mockAuthService.On("VerifyPhoneNumberAuthenticationToken", "testtoken1234").Return(true, nil)
		mockAuthService.On("ResetPassword", requestBody).Return(nil)

		app.Post(path, ResetPasswordHandler(mockAuthService))

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer testtoken1234")

		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	})

	t.Run("should return 400 when authorization header is not valid ", func(t *testing.T) {
		app := fiber.New()
		mockAuthService := new(mocks.MockAuthService)
		requestBody := dto.PasswordResetRequest{
			PhoneNumber: "+251923001100",
			NewPassword: "New password",
		}

		mockAuthService.On("VerifyPhoneNumberAuthenticationToken", "testtoken1234").Return(false, nil)
		mockAuthService.On("ResetPassword", requestBody).Return(nil)

		app.Post(path, ResetPasswordHandler(mockAuthService))

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer testtoken1234")

		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	})

	t.Run("should return 500 there's a problem resetting the password", func(t *testing.T) {
		app := fiber.New()
		mockAuthService := new(mocks.MockAuthService)
		requestBody := dto.PasswordResetRequest{
			PhoneNumber: "+251923001100",
			NewPassword: "Password",
		}

		mockAuthService.On("VerifyPhoneNumberAuthenticationToken", "testtoken1234").Return(true, nil)
		mockAuthService.On("ResetPassword", requestBody).Return(app_errors.ErrRecordNotFound)

		app.Post(path, ResetPasswordHandler(mockAuthService))

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer testtoken1234")

		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	})

	t.Run("should return 200 when successful", func(t *testing.T) {
		app := fiber.New()
		mockAuthService := new(mocks.MockAuthService)
		requestBody := dto.PasswordResetRequest{
			PhoneNumber: "+251923001100",
			NewPassword: "Password",
		}

		mockAuthService.On("VerifyPhoneNumberAuthenticationToken", "testtoken1234").Return(true, nil)
		mockAuthService.On("ResetPassword", requestBody).Return(nil)

		app.Post(path, ResetPasswordHandler(mockAuthService))

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer testtoken1234")

		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	})

}

func TestChangePasswordHandler(t *testing.T) {
	path := "/change-password"

	t.Run("should return a 400 response when request body is invalid", func(t *testing.T) {

		app := fiber.New()
		mockAuthService := new(mocks.MockAuthService)
		requestBody := fiber.Map{}
		app.Post(path, ResetPasswordHandler(mockAuthService))

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))

		req.Header.Set("Content-Type", "application/json")

		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	})
	t.Run("should return a 400 reponse when changePasswordRequest is invalid", func(t *testing.T) {
		app := fiber.New()
		mockAuthService := new(mocks.MockAuthService)
		requestBody := dto.ChangePasswordRequest{
			NewPassword: "1234",
		}

		app.Post(path, ResetPasswordHandler(mockAuthService))

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))

		req.Header.Set("Content-Type", "application/json")

		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	})

	t.Run("should return a 400 response when password is invalid", func(t *testing.T) {
		user := entity.User{
			PhoneNumber: "+251923001100",
			Password:    "hashed password",
			FirstName:   "John",
			LastName:    "Doe",
			Role:        entity.Role{Name: "attendant"},
		}
		app := fiber.New()
		mockAuthService := new(mocks.MockAuthService)
		requestBody := dto.ChangePasswordRequest{
			NewPassword: "1234",
			OldPassword: "4321",
		}

		mockAuthService.On("ChangePassword", requestBody, user).Return(app_errors.ErrInvalidPassword)

		app.Post(path, func(c *fiber.Ctx) error {
			c.Locals("user", &user)
			result := ChangePasswordHandler(mockAuthService)(c)
			return result
		})

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))

		req.Header.Set("Content-Type", "application/json")

		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	})

	t.Run("should return a 500 response when an internal error occurs", func(t *testing.T) {
		user := entity.User{
			PhoneNumber: "+251923001100",
			Password:    "hashed password",
			FirstName:   "John",
			LastName:    "Doe",
			Role:        entity.Role{Name: "attendant"},
		}
		app := fiber.New()
		mockAuthService := new(mocks.MockAuthService)
		requestBody := dto.ChangePasswordRequest{
			NewPassword: "1234",
			OldPassword: "4321",
		}

		mockAuthService.On("ChangePassword", requestBody, user).Return(errors.New("internal error"))

		app.Post(path, func(c *fiber.Ctx) error {
			c.Locals("user", &user)
			result := ChangePasswordHandler(mockAuthService)(c)
			return result
		})

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))

		req.Header.Set("Content-Type", "application/json")

		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	})

	t.Run("should return a 200 response when successful", func(t *testing.T) {
		user := entity.User{
			PhoneNumber: "+251923001100",
			Password:    "hashed password",
			FirstName:   "John",
			LastName:    "Doe",
			Role:        entity.Role{Name: "attendant"},
		}
		app := fiber.New()
		mockAuthService := new(mocks.MockAuthService)
		requestBody := dto.ChangePasswordRequest{
			NewPassword: "1234",
			OldPassword: "4321",
		}

		mockAuthService.On("ChangePassword", requestBody, user).Return(nil)

		app.Post(path, func(c *fiber.Ctx) error {
			c.Locals("user", &user)
			result := ChangePasswordHandler(mockAuthService)(c)
			return result
		})

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))

		req.Header.Set("Content-Type", "application/json")

		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	})

}
