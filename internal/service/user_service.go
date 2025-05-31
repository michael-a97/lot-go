package service

import (
	"errors"
	"lot/api/dto"
	"lot/internal/entity"
	"lot/internal/errors"
	"lot/internal/repository"
	"lot/internal/utilities"
)

type userService struct {
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
}

func NewUserService(
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
) UserService {
	return &userService{
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
}

type UserService interface {
	SignUp(request dto.SignUpRequest) (*dto.UserDto, error)
}

func (u userService) SignUp(request dto.SignUpRequest) (*dto.UserDto, error) {
	role, err := u.roleRepository.FindUserRoleByName(request.Role)
	if err != nil {
		return nil, errors.New("role not found")
	}

	if _, err := u.userRepository.FindByPhoneNumber(request.PhoneNumber); !errors.Is(err, app_errors.ErrRecordNotFound) {
		return nil, errors.New("an account with that phone number already exists")
	}

	hashedPassword, err := utilities.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		Password:    hashedPassword,
		PhoneNumber: request.PhoneNumber,
		Role:        *role,
		RoleID:      role.ID,
	}

	savedUser, err := u.userRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return toDto(*savedUser), nil
}

func toDto(user entity.User) *dto.UserDto {
	return &dto.UserDto{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role.Name,
		ID:          user.ID,
	}
}
