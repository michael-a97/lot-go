package repository

import (
	"lot/internal/entity"

	"gorm.io/gorm"
)

type authRepository struct {
	DB *gorm.DB
}


type AuthRepository interface {
	RevokeAllRefreshTokensForUser(userId uint) error
	IsValidRefreshToken(token string, userId uint) (bool, error)
	SaveRefreshToken(token string, userId uint) error
}

func (a authRepository) SaveRefreshToken(token string, userId uint) error {
	refreshToken := entity.RefreshToken{
		Token:   token,
		Revoked: false,
		UserID:  userId,
	}
	result := a.DB.Save(&refreshToken)
	if result.Error != nil {
		return result.Error
	}
	return nil 
}

func (a authRepository) RevokeAllRefreshTokensForUser(userId uint) error {
	result := a.DB.Model(&entity.RefreshToken{}).Where("user_id = ?", userId).Update("revoked", true)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a authRepository) IsValidRefreshToken(token string, userId uint) (bool, error) {
	var savedToken *entity.RefreshToken

	result := a.DB.Model(
		&entity.RefreshToken{},
	).Where(
		"user_id = ?", userId,
	).Where(
		"token = ?", token,
	).Where("revoked = ?", false).Find(&savedToken)

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 1 {
		return true, nil
	}

	return false, nil

}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return authRepository{DB: db}
}
