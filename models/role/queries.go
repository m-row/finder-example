package role

import (
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder"
	"github.com/m-row/finder-example/model"
)

var (
	selects = &[]string{
		"roles.*",
	}
	inserts = &[]string{
		"name",
	}
)

func buildInput(m *Model) (*[]any, error) {
	input := &[]any{
		m.Name,
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
	if err := finder.ShowOne(shown, c); err != nil {
		return err
	}
	return m.GetPermissions(shown)
}

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
	}
	if err := finder.CreateOne(created, c); err != nil {
		return err
	}
	return m.SyncPermissions(created, tx)
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
	if err := finder.UpdateOne(updated, c); err != nil {
		return err
	}
	if len(updated.Permissions) != 0 {
		if err := m.SyncPermissions(updated, tx); err != nil {
			return err
		}
	}
	return nil
}

func (m *Queries) DeleteOne(deleted *Model, tx finder.Connection) error {
	c := &finder.ConfigDelete{
		DB:      tx,
		QB:      m.QB,
		Selects: selects,
	}
	return finder.DeleteOne(deleted, c)
}
