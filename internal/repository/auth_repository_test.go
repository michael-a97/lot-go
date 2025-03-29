package repository

import (
	"lot/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"

)



func TestSaveRefreshToken(t *testing.T) {
	t.Run("Should return nill when successful", func(t *testing.T) {
		token := "test token 1"
		userId := uint(1)
		db := setupTestDB()
		authRepository := NewAuthRepository(db)

		err := authRepository.SaveRefreshToken(token, userId)

		assert.Nil(t, err)
	})

	t.Run("Should return an error when an error occurs", func(t *testing.T) {
		token := "test token 1"
		userId := uint(1)
		db := setupTestDB()
		authRepository := NewAuthRepository(db)
		sqlDb, _ := db.DB()
		sqlDb.Close()

		err := authRepository.SaveRefreshToken(token, userId)

		assert.Error(t, err)
	})
}

func TestRevokeAllRefreshTokensForUser(t *testing.T) {
	t.Run("Should revoke previously saved tokens ", func(t *testing.T) {
		db := setupTestDB()
		authRepository := authRepository{DB: db}
		authRepository.SaveRefreshToken("token1", 1)
		authRepository.SaveRefreshToken("token2", 1)

		authRepository.RevokeAllRefreshTokensForUser(1)

		var activeTokens []entity.RefreshToken
		db.Model(&entity.RefreshToken{}).Where("revoked = ?", false).Find(&activeTokens)
		assert.Equal(t, len(activeTokens), 0)
	})

	t.Run("Should return an error when an error occurs ", func(t *testing.T) {
		db := setupTestDB()
		authRepository := authRepository{DB: db}
		authRepository.SaveRefreshToken("token1", 1)
		authRepository.SaveRefreshToken("token2", 1)
		sqlDB, _ := db.DB()

		sqlDB.Close()
		err := authRepository.RevokeAllRefreshTokensForUser(1)

		assert.Error(t, err)
	})

}

func TestIsValidRefreshToken(t *testing.T) {
	token := entity.RefreshToken{
		Token:  "token 1",
		UserID: 1,
	}

	t.Run("Should return true if the token is valid", func(t *testing.T) {
		db := setupTestDB()
		authRepository := NewAuthRepository(db)
		authRepository.SaveRefreshToken(token.Token, token.UserID)

		result, err := authRepository.IsValidRefreshToken(token.Token, token.UserID)

		assert.Equal(t, result, true)
		assert.Nil(t, err)
	})

	t.Run("Should return false if the token is invalid", func(t *testing.T) {
		db := setupTestDB()
		authRepository := NewAuthRepository(db)

		result, err := authRepository.IsValidRefreshToken(token.Token, token.UserID)

		assert.Equal(t, result, false)
		assert.Nil(t, err)
	})
	t.Run("Should return and error when an error occurs", func(t *testing.T) {
		db := setupTestDB()
		authRepository := NewAuthRepository(db)
		sqlDb, _ := db.DB()
		sqlDb.Close()

		result, err := authRepository.IsValidRefreshToken(token.Token, token.UserID)

		assert.Equal(t, result, false)
		assert.Error(t, err)
	})
}
