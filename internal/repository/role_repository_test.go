package repository

import (
	"lot/internal/entity"
	"lot/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test db")
	}

	db.AutoMigrate(&entity.Role{})
	return db
}

func TestRoleRepository_Save(t *testing.T) {
	role := entity.Role{
		Name: "attendant",
	}

	t.Run("Should return nil after saving successfully", func(t *testing.T) {
		db := setupTestDB()
		roleRepository := NewRoleRepository(db)
		err := roleRepository.Save(role)

		assert.Nil(t, err)
	})

	t.Run("Should an error when an error occurs", func(t *testing.T) {
		db := setupTestDB()
		roleRepository := NewRoleRepository(db)
		sqlDb, _ := db.DB()
		sqlDb.Close()

		err := roleRepository.Save(role)

		assert.Error(t, err)
	})
}

func TestFindUserRoleByName(t *testing.T) {
	role := entity.Role{
		Name: "attendant",
	}

	t.Run("Should return nil after saving successfully", func(t *testing.T) {
		db := setupTestDB()
		roleRepository := NewRoleRepository(db)
		roleRepository.Save(role)

		savedRole, err := roleRepository.FindUserRoleByName(role.Name)

		assert.Nil(t, err)
		assert.Equal(t, role.Name, savedRole.Name)
	})

	t.Run("Should an error when an error occurs", func(t *testing.T) {
		db := setupTestDB()
		roleRepository := NewRoleRepository(db)

		sqlDb, _ := db.DB()
		sqlDb.Close()

		role, err := roleRepository.FindUserRoleByName("test")

		assert.Nil(t, role)
		assert.Error(t, app_errors.ErrRecordNotFound, err)
	})
}

func TestFind(t *testing.T) {
	role := entity.Role{
		Name: "attendant",
	}

	t.Run("Should return nil after saving successfully", func(t *testing.T) {
		db := setupTestDB()
		roleRepository := NewRoleRepository(db)
		roleRepository.Save(role)

		roles := roleRepository.Find()

		assert.Equal(t, 1, len(roles))
	})

}
