package category

import (
	"context"
	"fmt"

	"github.com/m-row/finder-example/model"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder"
)

var inserts = &[]string{
	"id",
	"name",
	"parent_id",
	"super_parent_id",
	"sort",
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
		category.Sort,
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
	IsAdmin          bool
	SortBeforeUpdate int
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
	if err := finder.CreateOne(created, c); err != nil {
		return err
	}
	return m.fixAndSortOnInsert(created, conn)
}

func (m *Queries) UpdateOne(
	updated *Model,
	ws *WhereScope,
	conn finder.Connection,
) error {
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
	if err := finder.UpdateOne(updated, c); err != nil {
		return err
	}
	if ws.SortBeforeUpdate < updated.Sort {
		if err := m.changeSortInUpdate(
			ws.SortBeforeUpdate+1,
			updated.Sort+1,
			updated.Parent.ID,
			updated.Depth,
			false,
		); err != nil {
			return err
		}
	}
	if ws.SortBeforeUpdate > updated.Sort {
		if err := m.changeSortInUpdate(
			updated.Sort,
			ws.SortBeforeUpdate,
			updated.Parent.ID,
			updated.Depth,
			true,
		); err != nil {
			return err
		}
	}
	return nil
}

func (m *Queries) DeleteOne(
	deleted *Model,
	ws *WhereScope,
	conn finder.Connection,
) error {
	c := &finder.ConfigDelete{
		DB:         m.DB,
		QB:         m.QB,
		TableAlias: "c1",
		Joins:      joins("c1"),
		Wheres:     wheres("c1", ws),
		Selects:    selects("c1"),
	}
	if err := finder.DeleteOne(deleted, c); err != nil {
		return err
	}
	if _, err := conn.ExecContext(
		context.Background(),
		`
           UPDATE categories 
              SET sort = sort - 1 
            WHERE parent_id = $1 
              AND depth = $2 
              AND sort > $3
        `,
		deleted.Parent.ID,
		deleted.Depth,
		deleted.Sort,
	); err != nil {
		return err
	}
	return nil
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

func (m *Queries) changeSortInUpdate(
	starter, length int,
	parentID *uuid.UUID,
	depth int,
	isIncrease bool,
) error {
	var valuesSorted []int
	for i := starter; i < length; i++ {
		valuesSorted = append(valuesSorted, i)
	}

	var expression string
	if isIncrease {
		expression = "sort + 1"
	} else {
		expression = "sort - 1"
	}

	query, args, err := m.QB.Update("categories").
		Set("sort", squirrel.Expr(expression)).
		Where(squirrel.Eq{"sort": valuesSorted}).
		Where("parent_id = ?", parentID).
		Where("depth = ?", depth).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := m.DB.ExecContext(
		context.Background(),
		query,
		args...,
	); err != nil {
		return err
	}

	return nil
}

func (m *Queries) fixAndSortOnInsert(
	created *Model,
	conn finder.Connection,
) error {
	var inSequence bool
	if err := conn.GetContext(
		context.Background(),
		&inSequence,
		`
			SELECT
                -- reducing 2 has to do with the newly inserted item
				COALESCE(count(*)-2 = MAX(sort), false) AS inSequence
			FROM
				categories
			WHERE
				parent_id = $1 
                AND depth = $2
		`,
		created.Parent.ID,
		created.Depth,
	); err != nil {
		return err
	}

	if inSequence {
		if _, err := conn.ExecContext(
			context.Background(),
			`
               UPDATE categories 
                  SET sort = sort + 1 
                WHERE parent_id = $2 
                  AND depth = $3
                  AND id != $1
            `,
			created.ID,
			created.Parent.ID,
			created.Depth,
		); err != nil {
			return err
		}
		return nil
	}

	var ids []string
	if err := conn.SelectContext(
		context.Background(),
		&ids,
		`
            SELECT id 
            FROM categories 
            WHERE parent_id = $1 
                  AND depth = $2 
            ORDER BY sort
        `,
		created.Parent.ID,
		created.Depth,
	); err != nil {
		return err
	}
	queryUpdate := ""
	if len(ids) > 0 {
		for i, id := range ids {
			queryUpdate += fmt.Sprintf(
				`UPDATE categories SET sort = %d WHERE id = '%s';`,
				i,
				id,
			)
		}
		if _, err := conn.ExecContext(
			context.Background(),
			queryUpdate,
		); err != nil {
			return err
		}
	}
	return nil
}
