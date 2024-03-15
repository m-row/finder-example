package role

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/m-row/finder"
)

func (m *Queries) GrantAllPermissions(roleID int) (int64, error) {
	perms := []int{}

	if err := m.DB.SelectContext(
		context.Background(),
		&perms,
		`select id from permissions`,
	); err != nil {
		return 0, err
	}

	inserts := m.QB.
		Insert("role_permissions").
		Columns("role_id", "permission_id")
	for _, v := range perms {
		inserts = inserts.Values(roleID, v)
	}
	inserts = inserts.Suffix(`ON CONFLICT DO NOTHING`)
	query, args, err := inserts.ToSql()
	if err != nil {
		return 0, err
	}
	result, err := m.DB.ExecContext(context.Background(), query, args...)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (m *Queries) GrantByScope(roleID int, scope string) (int64, error) {
	perms := []int{}

	if err := m.DB.SelectContext(
		context.Background(),
		&perms,
		`
            SELECT id FROM permissions WHERE scope = $1
        `,
		scope,
	); err != nil {
		return 0, err
	}

	inserts := m.QB.
		Insert("role_permissions").
		Columns("role_id", "permission_id")
	for _, v := range perms {
		inserts = inserts.Values(roleID, v)
	}
	inserts = inserts.Suffix(`ON CONFLICT DO NOTHING`)
	query, args, err := inserts.ToSql()
	if err != nil {
		return 0, err
	}
	result, err := m.DB.ExecContext(context.Background(), query, args...)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (m *Queries) GetPermissions(role *Model) error {
	role.Permissions = []int{}

	query := `
        SELECT
            permissions.id
        FROM
            permissions
            INNER JOIN role_permissions r ON r.permission_id = permissions.id 
            INNER JOIN roles ON r.role_id = roles.id 
        WHERE
            roles.id = $1
        ORDER BY
            id
	`

	return m.DB.SelectContext(
		context.Background(),
		&role.Permissions,
		query,
		role.ID,
	)
}

func (m *Queries) SyncPermissions(role *Model, tx finder.Connection) error {
	if _, err := tx.ExecContext(
		context.Background(),
		`DELETE FROM role_permissions WHERE role_id = $1`,
		role.ID,
	); err != nil {
		return err
	}
	query, args, err := m.QB.
		Select("id").
		From("permissions").
		Where(squirrel.Eq{"id": role.Permissions}).
		ToSql()
	if err != nil {
		return err
	}
	var permissions []int
	if err := tx.SelectContext(
		context.Background(),
		&permissions,
		query,
		args...,
	); err != nil {
		return err
	}
	if len(permissions) == 0 {
		return nil
	}
	rolePerms := m.QB.
		Insert("role_permissions").
		Columns("role_id", "permission_id").
		Suffix("RETURNING permission_id as id")

	for _, value := range permissions {
		rolePerms = rolePerms.Values(role.ID, value)
	}

	query2, args2, err := rolePerms.ToSql()
	if err != nil {
		return err
	}
	if err := tx.SelectContext(
		context.Background(),
		&role.Permissions,
		query2,
		args2...,
	); err != nil {
		return err
	}
	return nil
}

func (m *Queries) RevokeAllPermissions(roleID int) (int64, error) {
	result, err := m.DB.ExecContext(
		context.Background(),
		`
            DELETE FROM role_permissions WHERE role_id = $1
        `,
		roleID,
	)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}
