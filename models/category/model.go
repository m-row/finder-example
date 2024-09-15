package category

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/m-row/finder-example/model"
	"github.com/m-row/finder-example/types"

	"github.com/google/uuid"
	"github.com/m-row/finder"
)

type Model struct {
	ID            uuid.UUID   `db:"id"              json:"id"`
	Name          types.JSONB `db:"name"            json:"name"`
	Depth         int         `db:"depth"           json:"depth"`
	IsDisabled    bool        `db:"is_disabled"     json:"is_disabled"`
	IsFeatured    bool        `db:"is_featured"     json:"is_featured"`
	ParentID      *uuid.UUID  `db:"parent_id"       json:"-"`
	SuperParentID *uuid.UUID  `db:"super_parent_id" json:"-"`
	CreatedAt     time.Time   `db:"created_at"      json:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at"      json:"updated_at"`

	Parent      MinimalModel `db:"parent"       json:"parent"`
	SuperParent MinimalModel `db:"super_parent" json:"super_parent"`
}

type MinimalModel struct {
	ID   *uuid.UUID   `db:"id"   json:"id"`
	Name *types.JSONB `db:"name" json:"name"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) ModelName() string {
	return "category"
}

func (m *Model) TableName() string {
	return "categories"
}

func (m *Model) DefaultSearch() string {
	return "name->>'ar'"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		m.DefaultSearch(),
		"name->>'en'",
	}
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{}
}

func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	isInsert := m.CreatedAt.Equal(time.Time{})
	if isInsert && m.ID == uuid.Nil {
		model.InputOrNewUUID(&m.ID, v)
	}
	return isInsert
}

// AssignSuperParent gets parent super_parent and depth assigned to body.
func (m *Model) AssignSuperParent(conn finder.Connection) error {
	if m.Parent.ID != nil {
		var parent Model
		if err := conn.GetContext(
			context.Background(),
			&parent,
			`
                SELECT 
                    id,
                    name,
                    parent_id,
                    super_parent_id,
                    depth
                FROM 
                    categories 
                WHERE 
                    id=$1
            `,
			m.Parent.ID,
		); err != nil {
			return err
		}
		if parent.Depth == 0 {
			m.SuperParent.ID = &parent.ID
		} else {
			m.SuperParent.ID = parent.SuperParentID
		}
		if m.Parent.ID != nil {
			if *m.Parent.ID == m.ID {
				return errors.New("category can't be a parent to itself")
			}
		}
		m.Depth = parent.Depth + 1
	}
	return nil
}
