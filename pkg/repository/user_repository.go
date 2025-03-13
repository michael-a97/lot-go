package repository

import (
	"errors"
	"gorm.io/gorm"
	"lot/pkg/entity"
)

type userRepository struct {
	DB *gorm.DB
}

type UserRepository interface {
	Save(user entity.User) error
	Delete(user entity.User) error
	FindById(id uint) (*entity.User, error)
	FindByPhoneNumber(phoneNumber string) *entity.User
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

var errRecordNotFound = errors.New("record not found")

func (u userRepository) FindById(id uint) (*entity.User, error) {
	var user entity.User
	result := u.DB.Model(&entity.User{}).Preload("Role").Find(&user, id)

	if result.RowsAffected == 0 {
		return nil, errRecordNotFound
	}
	return &user, nil
}

func (u userRepository) Save(user entity.User) error {
	result := u.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u userRepository) Delete(user entity.User) error {
	result := u.DB.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u userRepository) FindByPhoneNumber(phoneNumber string) *entity.User {
	var user entity.User
	result := u.DB.Model(&entity.User{}).Where("phone_number = ?", phoneNumber).Preload("Role").Find(&user)
	if result.RowsAffected == 0 {
		return nil
	} else {
		return &user
	}
}
