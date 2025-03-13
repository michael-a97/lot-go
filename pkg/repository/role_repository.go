package repository

import (
	entity "lot/pkg/entity"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Save(role entity.Role) error
	Find() []entity.Role
	FindUserRoleByName(name string) (entity.Role, error)

}

type roleRepository struct {
	DB *gorm.DB
}

func (r roleRepository) Save(role entity.Role) error {
	err := r.DB.Save(&role)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r roleRepository) Find() []entity.Role {
	var roles []entity.Role
	r.DB.Model(&entity.Role{}).Find(&roles)
	return roles
}

func (r roleRepository) FindUserRoleByName(name string) (entity.Role, error) {
	var role entity.Role
	result := r.DB.Where("name = ?", name).First(&role)
	if result.RowsAffected == 0 {
		return entity.Role{}, errRecordNotFound
	}
	return role, nil
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return roleRepository{DB: db}
}
