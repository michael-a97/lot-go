package repository

import "github.com/stretchr/testify/mock"

type MockAuthRepository struct {
	mock.Mock
}

func (mockAuthRepository *MockAuthRepository) SaveRefreshToken(token string, userId uint) error {
	args := mockAuthRepository.Called(token, userId)
	return args.Error(0)
}

func (mockAuthRepository *MockAuthRepository) RevokeAllRefreshTokensForUser(userId uint) error {
	args := mockAuthRepository.Called(userId)
	return args.Error(0)
}

func (mockAuthRepository *MockAuthRepository) IsValidRefreshToken(token string, userId uint) (bool, error) {
	args := mockAuthRepository.Called(userId)
	return args.Bool(0), args.Error(1)
}
