package category

import (
	"context"

	"github.com/m-row/finder-example/model"

	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder"
)

var inserts = &[]string{
	"id",
	"name",
	"parent_id",
	"super_parent_id",
	"depth",
	"is_disabled",
	"is_featured",
}

func buildInput(category *Model) (*[]any, error) {
	input := &[]any{
		category.ID,
		category.Name,
		category.Parent.ID,
		category.SuperParent.ID,
		category.Depth,
		category.IsDisabled,
		category.IsFeatured,
	}
	if len(*input) != len(*inserts) {
		return nil, finder.ErrInputLengthMismatch(input, inserts)
	}
	return input, nil
}

func joins(alias string) *[]string {
	return &[]string{
		"categories as p ON " + alias + ".parent_id = p.id",
		"categories as sp ON " + alias + ".super_parent_id = sp.id",
	}
}

func selects(alias string) *[]string {
	return &[]string{
		alias + ".*",

		"p.id as \"parent.id\"",
		"p.name as \"parent.name\"",

		"sp.id as \"super_parent.id\"",
		"sp.name as \"super_parent.name\"",
	}
}

type WhereScope struct {
	IsAdmin bool
}

func wheres(alias string, ws *WhereScope) *[]squirrel.Sqlizer {
	w := []squirrel.Sqlizer{}
	if ws == nil {
		return &w
	}
	if ws.IsAdmin {
		return &w
	}
	w = append(w, squirrel.Expr(alias+".is_disabled=false"))
	return &w
}

type Queries struct {
	*model.Dependencies
}

func New(d *model.Dependencies) *Queries {
	return &Queries{d}
}

func (m *Queries) GetAll(
	ctx echo.Context,
	ws *WhereScope,
) (*finder.IndexResponse[*Model], error) {
	cfg := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		PGInfo:  m.PGInfo,
		Joins:   joins("categories"),
		Wheres:  wheres("categories", ws),
		Selects: selects("categories"),
		GroupBys: &[]string{
			"categories.id",
			"p.id",
			"sp.id",
		},
	}
	return finder.IndexBuilder[*Model](ctx.QueryParams(), cfg)
}

func (m *Queries) GetOne(
	shown *Model,
	ws *WhereScope,
) error {
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Joins:   joins("categories"),
		Wheres:  wheres("categories", ws),
		Selects: selects("categories"),
	}
	return finder.ShowOne(shown, c)
}

func (m *Queries) CreateOne(created *Model, conn finder.Connection) error {
	if err := created.AssignSuperParent(conn); err != nil {
		return err
	}
	input, err := buildInput(created)
	if err != nil {
		return err
	}
	c := &finder.ConfigStore{
		DB:         conn,
		QB:         m.QB,
		Input:      input,
		Inserts:    inserts,
		TableAlias: "c1",
		Joins:      joins("c1"),
		Selects:    selects("c1"),
	}
	return finder.CreateOne(created, c)
}

func (m *Queries) UpdateOne(
	updated *Model,
	ws *WhereScope,
	conn finder.Connection,
) error {
	if err := updated.AssignSuperParent(conn); err != nil {
		return err
	}
	input, err := buildInput(updated)
	if err != nil {
		return err
	}
	c := &finder.ConfigUpdate{
		DB:         conn,
		QB:         m.QB,
		Input:      input,
		Inserts:    inserts,
		TableAlias: "c1",
		Joins:      joins("c1"),
		Wheres:     wheres("c1", ws),
		Selects:    selects("c1"),
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	return finder.UpdateOne(updated, c)
}

func (m *Queries) DeleteOne(
	deleted *Model,
	ws *WhereScope,
	conn finder.Connection,
) error {
	c := &finder.ConfigDelete{
		DB:         conn,
		QB:         m.QB,
		TableAlias: "c1",
		Joins:      joins("c1"),
		Wheres:     wheres("c1", ws),
		Selects:    selects("c1"),
	}
	return finder.DeleteOne(deleted, c)
}

func (m *Queries) HasChildren(category *Model) (bool, error) {
	count := 0
	query, args, err := m.QB.
		Select("COUNT(*)").
		From("categories").
		Where("parent_id = ?", category.ID).
		ToSql()
	if err != nil {
		return false, err
	}

	if err := m.DB.GetContext(
		context.Background(),
		&count,
		query,
		args...,
	); err != nil {
		return count > 0, err
	}
	return count > 0, nil
}
