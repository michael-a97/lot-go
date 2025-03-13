package service

import (
	entity "lot/pkg/entity"
	repository "lot/pkg/repository"
)

type RoleService interface {
	Save(role entity.Role) error
	Find() []entity.Role
}

type roleService struct {
	RoleRepository repository.RoleRepository
}

func (r roleService) Save(role entity.Role) error {
	return r.RoleRepository.Save(role)
}

func (r roleService) Find() []entity.Role {
	return r.RoleRepository.Find()
}

func NewRoleService(roleRepository repository.RoleRepository) RoleService {
	return roleService{
		RoleRepository: roleRepository,
	}
}
