package service

import "github.com/stretchr/testify/mock"

type MockSmsTokenVerifier struct {
	mock.Mock
}

func (m *MockSmsTokenVerifier) IsTokenValid(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}
