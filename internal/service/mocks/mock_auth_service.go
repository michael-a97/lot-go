package mocks

import (
	"lot/api/dto"
	"lot/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) SignIn(request dto.LoginRequest) (*dto.AuthenticationResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*dto.AuthenticationResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(request dto.TokenRefreshRequest) (*dto.AuthenticationResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*dto.AuthenticationResponse), args.Error(1)
}

func (m *MockAuthService) VerifyPhoneNumberAuthenticationToken(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) ResetPassword(request dto.PasswordResetRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockAuthService) ChangePassword(request dto.ChangePasswordRequest, user entity.User) error {
	args := m.Called(request, user)
	return args.Error(0)
}

func (m *MockAuthService) GetUserFromAccessToken(accessToken string) (*entity.User, error) {
	args := m.Called(accessToken)
	return args.Get(0).(*entity.User), args.Error(1)
}
