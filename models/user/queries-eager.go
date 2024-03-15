package user

import (
	"context"

	"github.com/Masterminds/squirrel"
)

func (m *Queries) EagerLoad(
	list *[]*Model,
) error {
	var ids []string
	for _, v := range *list {
		ids = append(ids, v.ID.String())
	}
	return m.GetRolesForList(list, ids)
}

func (m *Queries) GetRolesForList(
	list *[]*Model,
	ids []string,
) error {
	query, args, err := m.QB.
		Select(
			"user_id",
			"role_id",
		).
		From("user_roles").
		Where(squirrel.Eq{"user_id": ids}).
		ToSql()
	if err != nil {
		return err
	}

	userRoles := []UserRole{}

	if err := m.DB.SelectContext(
		context.Background(),
		&userRoles,
		query,
		args...,
	); err != nil {
		return err
	}
	for i := range *list {
		v := (*list)[i]
		v.Roles = &[]int{}
		for j := range userRoles {
			ur := userRoles[j]
			if ur.UserID == v.ID {
				*v.Roles = append(*v.Roles, ur.RoleID)
			}
		}
	}
	return nil
}
