package handler

import (
	"bytes"
	"encoding/json"
	"lot/api/dto"
	app_errors "lot/internal/errors"
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
