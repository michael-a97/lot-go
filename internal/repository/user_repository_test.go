package repository

import (
	"lot/internal/entity"
	"lot/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDatabase_UserRepository() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test db")
	}

	_ = db.AutoMigrate(&entity.User{})
	return db
}

func TestSave(t *testing.T) {
	role := entity.Role{
		Name: "attendant",
	}
	user1 := entity.User{
		PhoneNumber: "+251923001100",
		Password:    "hashed password",
		FirstName:   "John",
		LastName:    "Doe",
		Role:        role,
	}

	t.Run("Should return nil after saving successfully", func(t *testing.T) {
		db := setupTestDatabase_UserRepository()
		userRepository := NewUserRepository(db)
		err := userRepository.Save(user1)

		assert.Nil(t, err)
	})

	t.Run("Should an error when an error occurs", func(t *testing.T) {
		db := setupTestDatabase_UserRepository()
		userRepository := NewUserRepository(db)
		sqlDb, _ := db.DB()
		_ = sqlDb.Close()

		err := userRepository.Save(user1)

		assert.Error(t, err)
	})
}

func TestFindById(t *testing.T) {
	role := entity.Role{
		Name: "attendant",
	}
	user1 := entity.User{
		PhoneNumber: "+251923001100",
		Password:    "hashed password",
		FirstName:   "John",
		LastName:    "Doe",
		Role:        role,
	}

	t.Run("Should return previously saved user", func(t *testing.T) {
		db := setupTestDatabase_UserRepository()
		userRepository := NewUserRepository(db)
		_ = userRepository.Save(user1)

		user, err := userRepository.FindById(1)

		assert.Equal(t, user1.FirstName, user.FirstName)
		assert.Equal(t, user1.LastName, user.LastName)
		assert.Equal(t, user1.PhoneNumber, user.PhoneNumber)
		assert.Equal(t, user1.Password, user.Password)
		assert.Equal(t, user1.Role.Name, user.Role.Name)
		assert.Nil(t, err)
	})

	t.Run("Should return nil, record not found error when there are no results found", func(t *testing.T) {
		db := setupTestDatabase_UserRepository()
		userRepository := NewUserRepository(db)

		user, err := userRepository.FindById(1)

		assert.Nil(t, user)
		assert.Equal(t, app_errors.ErrRecordNotFound, err)
	})
}

func TestDelete(t *testing.T) {
	role := entity.Role{
		Name: "attendant",
	}
	user1 := entity.User{
		PhoneNumber: "+251923001100",
		Password:    "hashed password",
		FirstName:   "John",
		LastName:    "Doe",
		Role:        role,
	}

	t.Run("Should return nil after deleting successfully", func(t *testing.T) {
		db := setupTestDatabase_UserRepository()
		userRepository := NewUserRepository(db)
		_ = userRepository.Save(user1)
		savedUser, _ := userRepository.FindById(1)

		err := userRepository.Delete(*savedUser)

		assert.Nil(t, err)
	})

	t.Run("Should an error when an error occurs", func(t *testing.T) {
		db := setupTestDatabase_UserRepository()
		userRepository := NewUserRepository(db)
		_ = userRepository.Save(user1)
		savedUser, _ := userRepository.FindById(1)
		sqlDb, _ := db.DB()
		_ = sqlDb.Close()

		err := userRepository.Delete(*savedUser)

		assert.Error(t, err)
	})
}

func TestFindByPhoneNumber(t *testing.T) {
	role := entity.Role{
		Name: "attendant",
	}
	user1 := entity.User{
		PhoneNumber: "+251923001100",
		Password:    "hashed password",
		FirstName:   "John",
		LastName:    "Doe",
		Role:        role,
	}

	t.Run("Should return previously saved user", func(t *testing.T) {
		db := setupTestDatabase_UserRepository()
		userRepository := NewUserRepository(db)
		_ = userRepository.Save(user1)

		user, err := userRepository.FindByPhoneNumber(user1.PhoneNumber)

		assert.Equal(t, user1.FirstName, user.FirstName)
		assert.Equal(t, user1.LastName, user.LastName)
		assert.Equal(t, user1.PhoneNumber, user.PhoneNumber)
		assert.Equal(t, user1.Password, user.Password)
		assert.Equal(t, user1.Role.Name, user.Role.Name)
		assert.Nil(t, err)
	})

	t.Run("Should return nil, record not found error when there are no results found", func(t *testing.T) {
		db := setupTestDatabase_UserRepository()
		userRepository := NewUserRepository(db)

		user, err := userRepository.FindByPhoneNumber(user1.PhoneNumber)

		assert.Nil(t, user)
		assert.Equal(t, app_errors.ErrRecordNotFound, err)
	})
}
