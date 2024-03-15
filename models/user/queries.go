package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder"
	"github.com/m-row/finder-example/model"
)

var (
	selects = &[]string{
		"users.id",
		"users.name",
		"users.phone",
		"users.email",
		"users.password_hash",
		"users.is_disabled",
		"users.created_at",
		"users.updated_at",

		// TODO:
		// config.SQLSelectURLPath("users", "img", "img"),
		// config.SQLSelectURLPath("users", "thumb", "thumb"),
	}
	joins   = &[]string{}
	inserts = &[]string{
		"id",
		"email",
		"password_hash",
		"name",
		"phone",
		"is_disabled",
		"img",
		"thumb",
	}
)

func buildInput(m *Model) (*[]any, error) {
	hash := ""
	if m.Password.Hash != nil {
		// this sets the password hash on create/ or when provided in update
		hash = string(*m.Password.Hash)
	} else if m.PasswordHash != nil {
		// this handle updates when password is not updated
		hash = string(*m.PasswordHash)
	}
	input := &[]any{
		m.ID,
		squirrel.Expr("lower(?)", m.Email),
		hash,
		m.Name,
		m.Phone,
		m.IsDisabled,
		m.Img,
		m.Thumb,
	}
	if len(*input) != len(*inserts) {
		return nil, finder.ErrInputLengthMismatch(input, inserts)
	}
	return input, nil
}

type Queries struct {
	*model.Dependencies
}

func New(d *model.Dependencies) *Queries {
	return &Queries{d}
}

func (m *Queries) GetAll(
	ctx echo.Context,
) (*finder.IndexResponse[*Model], error) {
	cfg := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		PGInfo:  m.PGInfo,
		Joins:   joins,
		Selects: selects,
		GroupBys: &[]string{
			"users.id",
			// "wallets.id",
		},
	}
	indexResponse, err := finder.IndexBuilder[*Model](ctx.QueryParams(), cfg)
	if err != nil {
		return nil, err
	}
	if err := m.EagerLoad(indexResponse.Data); err != nil {
		return nil, err
	}
	return indexResponse, nil
}

func (m *Queries) GetOne(shown *Model) error {
	if shown.ID == uuid.Nil && shown.Phone == nil && shown.Email == nil {
		return nil
	}
	wheres := &[]squirrel.Sqlizer{}
	if shown.ID != uuid.Nil {
		*wheres = append(
			*wheres,
			squirrel.Expr("users.id=?", shown.ID.String()),
		)
	}
	if shown.Phone != nil {
		if *shown.Phone != "" {
			expr := squirrel.Expr("users.phone = ?", *shown.Phone)
			*wheres = append(*wheres, expr)
		}
	}
	if shown.Email != nil {
		if *shown.Email != "" {
			expr := squirrel.Expr("users.email = LOWER(?)", *shown.Email)
			*wheres = append(*wheres, expr)
		}
	}
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Joins:   joins,
		Wheres:  wheres,
		Selects: selects,
	}
	if err := finder.ShowOne(shown, c); err != nil {
		return err
	}
	return m.GetRoles(shown)
}

// CreateOne inserts a user with roles,
// db will create a cart and a wallet using a trigger
//
// requires a transaction.
func (m *Queries) CreateOne(created *Model, tx finder.Connection) error {
	input, err := buildInput(created)
	if err != nil {
		return err
	}
	c := &finder.ConfigStore{
		DB:      m.DB,
		QB:      m.QB,
		Input:   input,
		Inserts: inserts,
		Selects: selects,
		Joins:   joins,
	}
	if err := finder.CreateOne(created, c); err != nil {
		return err
	}
	return m.AssignRoles(created, tx)
}

func (m *Queries) UpdateOne(updated *Model, tx finder.Connection) error {
	input, err := buildInput(updated)
	if err != nil {
		return err
	}
	c := &finder.ConfigUpdate{
		DB:      tx,
		QB:      m.QB,
		Input:   input,
		Inserts: inserts,
		Selects: selects,
		Joins:   joins,
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	if err := finder.UpdateOne(updated, c); err != nil {
		return err
	}
	if updated.Roles != nil {
		if len(*updated.Roles) != 0 {
			if err := m.AssignRoles(updated, tx); err != nil {
				return err
			}
			if err := m.GetPermissions(updated, tx); err != nil {
				return err
			}
		}
	}
	return nil
}

// ClearOne Clear user account data.
func (m *Queries) ClearOne(userID *uuid.UUID, tx finder.Connection) error {
	// clear user data
	query1 := `
       UPDATE
           users
       SET
           phone = NULL,
           email = NULL,
           password_hash = NULL,
           name = NULL,
           is_disabled = TRUE,
           img = NULL,
           thumb = NULL
       WHERE
           id = $1
    `
	if _, err := tx.ExecContext(
		context.Background(),
		query1,
		userID,
	); err != nil {
		return err
	}
	// clear user roles
	if _, err := tx.ExecContext(
		context.Background(),
		`DELETE FROM user_roles WHERE user_id=$1`,
		*userID,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	// clear user permissions
	if _, err := tx.ExecContext(
		context.Background(),
		`DELETE FROM user_permissions WHERE user_id=$1`,
		*userID,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	return nil
}
