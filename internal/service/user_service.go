package service

import (
	"errors"
	"lot/api/dto"
	"lot/internal/entity"
	app_errors "lot/internal/errors"
	"lot/internal/repository"

	"golang.org/x/crypto/bcrypt"
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
	hashPassword(password string) (string, error)
}

func (u userService) SignUp(request dto.SignUpRequest) (*dto.UserDto, error) {
	role, err := u.roleRepository.FindUserRoleByName(request.Role)
	if err != nil {
		return nil, errors.New("role not found")
	}

	if _, err := u.userRepository.FindByPhoneNumber(request.PhoneNumber); err == app_errors.ErrRecordNotFound {
		return nil, errors.New("an account with that phone number already exists")
	}

	hashedPassword, err := u.hashPassword(request.Password)
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

	if err := u.userRepository.Save(user); err != nil {
		return nil, err
	}

	return toDto(user), nil
}

func (u userService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func toDto(user entity.User) *dto.UserDto {
	return &dto.UserDto{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role.Name,
	}
}
