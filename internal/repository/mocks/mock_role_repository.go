package repository

import (
	entity "lot/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Save(role entity.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) Find() []entity.Role {
	args := m.Called()
	return args.Get(0).([]entity.Role)
}

func (m *MockRoleRepository) FindUserRoleByName(name string) (*entity.Role, error) {
	args := m.Called()
	return args.Get(0).(*entity.Role), args.Error(1)
}

