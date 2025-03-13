package database

import (
	"fmt"
	"lot/pkg/entity"
	"lot/pkg/service"
)

func Seed(roleService service.RoleService) {
	existingRoles := roleService.Find()
	if len(existingRoles) == 0 {
		roles := [3]entity.Role{
			{Name: entity.RoleAdmin},
			{Name: entity.RoleAttendant},
			{Name: entity.RoleUser},
		}

		for _, role := range roles {
			err := roleService.Save(role)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

}
