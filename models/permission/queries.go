package permission

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder"
	"github.com/m-row/finder-example/model"
)

var (
	selects = &[]string{
		"permissions.*",
	}
	inserts = &[]string{
		"method",
		"path",
		"model",
		"action",
		"scope",
	}
)

func buildInput(m *Model) (*[]any, error) {
	input := &[]any{
		squirrel.Expr("UPPER(?)", m.Method),
		squirrel.Expr("LOWER(?)", m.Path),
		squirrel.Expr("LOWER(?)", m.Model),
		squirrel.Expr("LOWER(?)", m.Action),
		squirrel.Expr("LOWER(?)", m.Scope),
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
	config := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		PGInfo:  m.PGInfo,
		Selects: selects,
	}
	return finder.IndexBuilder[*Model](ctx.QueryParams(), config)
}

func (m *Queries) GetOne(shown *Model) error {
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Selects: selects,
	}
	return finder.ShowOne(shown, c)
}

func (m *Queries) CreateOne(created *Model, tx finder.Connection) error {
	input, err := buildInput(created)
	if err != nil {
		return err
	}
	c := &finder.ConfigStore{
		DB:      tx,
		QB:      m.QB,
		Input:   input,
		Inserts: inserts,
		Selects: selects,
	}
	return finder.CreateOne(created, c)
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
	}
	return finder.UpdateOne(updated, c)
}

func (m *Queries) DeleteOne(deleted *Model, tx finder.Connection) error {
	c := &finder.ConfigDelete{
		DB:      tx,
		QB:      m.QB,
		Selects: selects,
	}
	return finder.DeleteOne(deleted, c)
}

func (m *Queries) GetByList(ids *[]int) (*[]Model, error) {
	perms := []Model{}
	query, args, err := m.QB.
		Select(*selects...).
		From("permissions").
		Where(squirrel.Eq{"id": ids}).
		OrderBy("model", "action", "scope").
		ToSql()
	if err != nil {
		return nil, err
	}
	if err := m.DB.SelectContext(
		context.Background(),
		&perms,
		query,
		args...,
	); err != nil {
		return nil, err
	}
	return &perms, nil
}

func (m *Queries) GetByMethodPathList(
	method, path string,
	ids []int,
) (*[]Model, error) {
	perms := []Model{}
	if len(ids) == 0 {
		return &perms, nil
	}
	query, args, err := m.QB.
		Select(*selects...).
		From("permissions").
		Where("path = ?", path).
		Where("method = ?", method).
		Where(squirrel.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return nil, err
	}
	if err := m.DB.SelectContext(
		context.Background(),
		&perms,
		query,
		args...,
	); err != nil {
		return nil, err
	}
	return &perms, nil
}

func (m *Queries) BulkCreate(perms *[]Model) (int64, error) {
	inserts := m.QB.
		Insert("permissions").
		Columns(
			"method",
			"path",
			"model",
			"action",
			"scope",
		).
		Suffix(`ON CONFLICT DO NOTHING`)

		//  NOTE: can updated on conflict using this
		//   ON CONFLICT (name)
		//   DO UPDATE SET
		//       method = excluded.method,
		//       path = excluded.path,
		//       name = excluded.name,
		//       scope = excluded.scope

	if perms != nil {
		for _, v := range *perms {
			inserts = inserts.Values(
				v.Method,
				v.Path,
				v.Model,
				v.Action,
				v.Scope,
			)
		}
	}

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

func (m *Queries) DeleteUnused(perms *[]Model) (int64, error) {
	dbPerms := []Model{}

	if err := m.DB.SelectContext(
		context.Background(),
		&dbPerms,
		`SELECT * FROM permissions`,
	); err != nil {
		return 0, err
	}
	if perms == nil {
		return 0, errors.New("perms must be pointer to []Model")
	}

	discarded := []int{}
	for _, dp := range dbPerms {
		matches := 0
		for _, p := range *perms {
			t1 := dp.Model == p.Model
			t2 := dp.Action == p.Action
			t3 := dp.Path == p.Path
			t4 := dp.Method == p.Method
			t5 := dp.Scope == p.Scope

			if t1 && t2 && t3 && t4 && t5 {
				matches += 1
			}
		}
		if matches == 0 {
			discarded = append(discarded, dp.ID)
		}
	}

	query, args, err := m.QB.
		Delete("permissions").
		Where(squirrel.Eq{"id": discarded}).
		ToSql()
	if err != nil {
		return 0, err
	}
	result, err := m.DB.ExecContext(context.Background(), query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *Queries) DistinctScopes(scopes *[]string) error {
	return m.DB.SelectContext(
		context.Background(),
		scopes,
		`SELECT DISTINCT(scope) FROM permissions`,
	)
}
