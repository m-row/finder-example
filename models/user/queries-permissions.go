package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-row/finder"
)

func (m *Queries) GetPermissions(user *Model, tx finder.Connection) error {
	user.Permissions = &[]int{}
	query := `
    SELECT
        permissions.id
    FROM
        permissions
    INNER JOIN user_permissions up ON up.permission_id = permissions.id
    INNER JOIN users ON up.user_id = users.id
    UNION
        SELECT
            permissions.id AS permission_id
        FROM
            permissions
            INNER JOIN role_permissions r ON r.permission_id = permissions.id
            INNER JOIN roles ON r.role_id = roles.id
            INNER JOIN user_roles ON roles.id = user_roles.role_id
            INNER JOIN users ON user_roles.user_id = users.id
        WHERE
            users.id = $1
        ORDER BY
            id
	`

	conn := m.DB
	if tx != nil {
		conn = tx
	}

	return conn.SelectContext(
		context.Background(),
		user.Permissions,
		query,
		user.ID,
	)
}

func (m *Queries) RevokeAllPermissions(userID *uuid.UUID, tx *sqlx.Tx) error {
	if _, err := tx.ExecContext(
		context.Background(),
		`DELETE FROM user_permissions WHERE user_id = $1`,
		userID,
	); err != nil {
		return err
	}
	return nil
}

func (m *Queries) SyncPermissions(u *Model, tx *sqlx.Tx) error {
	if err := m.RevokeAllPermissions(&u.ID, tx); err != nil {
		return err
	}
	query, args, err := m.QB.
		Select("id").
		From("permissions").
		Where(squirrel.Eq{"id": u.Permissions}).
		ToSql()
	if err != nil {
		return err
	}
	var permissions []int
	if err := tx.Select(&permissions, query, args...); err != nil {
		return err
	} else {
		userPerms := m.QB.
			Insert("user_permissions").
			Columns("user_id", "permission_id")
		for _, value := range permissions {
			userPerms = userPerms.Values(u.ID, value)
		}

		query, args, err := userPerms.ToSql()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(
			context.Background(),
			query,
			args...,
		); err != nil {
			return err
		}
		return nil
	}
}

func (m *Queries) GrantAllPermissions(userID *uuid.UUID) (int64, error) {
	perms := []int{}

	if err := m.DB.SelectContext(
		context.Background(),
		&perms,
		`select id from permissions`,
	); err != nil {
		return 0, err
	}

	inserts := m.QB.
		Insert("user_permissions").
		Columns("user_id", "permission_id")
	for _, v := range perms {
		inserts = inserts.Values(userID, v)
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
	return result.RowsAffected()
}
